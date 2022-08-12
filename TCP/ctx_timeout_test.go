package TCP

import (
	"context"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// >> Creates a dialer with a timeout
// note this dialer will immidiately return an error due to the control function
func TimeoutDialer(network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

// >> tests the returned error from dialer timeout
func TestDialTimeout(t *testing.T) {
	duration := time.Second * 5
	c, err := TimeoutDialer("tcp", "172.217.204.102:1651", duration)

	if err == nil {
		c.Close()
		t.Fatal("Connected")
	}

	nErr, ok := err.(net.Error)
	require.True(t, ok)
	require.True(t, nErr.Timeout())
	require.True(t, nErr.Temporary())
}

// >> testing the dialer's context timeout
func TestDialContext(t *testing.T) {
	// >> creating deadline
	// For the context -> so when expires, context will cnacle
	expires := 2 * time.Second
	dl := time.Now().Add(expires)

	// >> Creating the context
	ctx, cancel := context.WithDeadline(context.Background(), dl)
	t.Log("context created")
	defer cancel() // good practice, to make sure ctx is garbage collected

	// >> Checks Context timeout
	go func() {
		<-ctx.Done()
		require.Error(t, ctx.Err())
		require.Error(t, ctx.Err(), context.DeadlineExceeded)
		t.Log("context expired")
	}()

	// >> Creating Dialer
	t.Log("creating dialer")
	var d net.Dialer
	// Sleep long enough to reach the context's deadline.
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		time.Sleep(expires + time.Millisecond)
		return nil
	}

	// >> Dialing with a context
	// the context is expired
	t.Log("dialing")
	conn, err := d.DialContext(ctx, "tcp", "127.0.0.1:")
	require.Error(t, err)
	require.Nil(t, conn)

	nErr, ok := err.(net.Error)
	require.True(t, ok)
	require.True(t, nErr.Timeout())

	// >> closing connection
	t.Log("closing connection")
}

// WIP
// >> Testing the cancelation of multiple client connections using a single context
func TestDialContextCancelFanOut(t *testing.T) {
	// >> creating deadline
	// For the context -> so when expires, context will cnacle
	expires := 2 * time.Second
	dl := time.Now().Add(expires)

	// >> creating context
	ctx, cancel := context.WithDeadline(context.Background(), dl)

	//>> Creating server
	t.Log("creating server")
	listener, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer listener.Close()

	// >> Accepting Connections & Closing it after success
	go func() {
		// Only accepting a single connection.
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	// >> creating Dialers
	dial := func(ctx context.Context, address string, response chan int,
		id int, wg *sync.WaitGroup) {

		defer wg.Done()
		var d net.Dialer

		c, err := d.DialContext(ctx, "tcp", address)
		//@ connection Fail:
		// Simply return
		//@ connection success:
		// Close connection, and place the Dialer id in the response channel
		if err != nil {
			return
		}
		c.Close()

		select {
		case <-ctx.Done():
		case response <- id:
		}
	}

	//>> Calling the Dialers
	res := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	//>> canceling ctx, and closing channel
	response := <-res
	cancel()
	wg.Wait()
	close(res)

	//>> checking err type
	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %s",
			ctx.Err(),
		)
	}

	t.Logf("dialer %d retrieved the resource", response)
}
