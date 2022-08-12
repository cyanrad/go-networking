package TCP

import (
	"context"
	"fmt"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

// >> write pings at regular intervals
// - ctx: for termination & leakage prevention
// - reset: to signal timer reset
func pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	// the interval time value
	var interval time.Duration

	// >> getting interval time
	// we put the initial interval duration in the reset channel
	select {
	case <-ctx.Done(): //terminating
		return
	case interval = <-reset: // pulled initial interval off reset channel
	}

	// >> interval duration handling
	if interval <= 0 {
		interval = defaultPingInterval
	}

	// >> creating ping timer
	timer := time.NewTimer(interval)
	defer func() { // drains timer channel to avoide leakage
		if !timer.Stop() {
			<-timer.C
		}
	}()

	// >> pinging loop
	// keep track of time-outs by passing the ctx's cancle func
	// and call it here if @ max concecutive timeouts
	for {
		select {
		case <-ctx.Done(): //terminate
			return
		case newInterval := <-reset: // resetting the timer (if data recieved)
			if !timer.Stop() {
				<-timer.C // Blocking wait until it finishes
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C: // ping (timer expires)
			if _, err := w.Write([]byte("ping")); err != nil {
				// track and act on consecutive timeouts here
				return
			}
		}
		_ = timer.Reset(interval) // reset timer
	}
}

// >> simple ping example
// use in main
func ExamplePinger() {
	// the reset ctx
	ctx, cancel := context.WithCancel(context.Background())

	//writer/ reader
	r, w := io.Pipe() // in lieu of net.Conn

	// >> channels
	done := make(chan struct{})
	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second // initial ping interval

	// >> gort pinging
	go func() {
		pinger(ctx, w, resetTimer)
		close(done)
	}()

	// >> reset timer & reading pings
	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			fmt.Printf("resetting timer (%s)\n", d)
			resetTimer <- d
		}
		now := time.Now()
		buf := make([]byte, 1024)
		n, err := r.Read(buf) // reading from io pipe, n is byte count
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("received %q (%s)\n",
			buf[:n], time.Since(now).Round(100*time.Millisecond))
	}
	for i, v := range []int64{0, 200, 300, 0, -1, -1, -1} {
		fmt.Printf("Run %d:\n", i+1)
		receivePing(time.Duration(v)*time.Millisecond, r)
	}
	cancel() //terminating pinger
	<-done   // ensures the pinger exits after canceling the context
	// Output:
	// Run 1:
	// resetting timer (0s)
	// received "ping" (1s)
	// Run 2:
	// resetting timer (200ms)
	// received "ping" (200ms)
	// Run 3:
	// resetting timer (300ms)
	// received "ping" (300ms)
	// Run 4:
	// resetting timer (0s)
	// received "ping" (300ms)
	// Run 5:
	// received "ping" (300ms)
	// Run 6:
	// received "ping" (300ms)
	// Run 7:
	// received "ping" (300ms)
}
