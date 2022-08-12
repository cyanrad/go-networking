package TCP

// >> Internet Control Message Protocol
// gives you feedback about NW condition.
// determine host connectivity by pinging them.
//
// most hosts filter/block ICMP, so
// it would return that the entire system is down.
// avoid this, is by establishing a TCP connection on some port.

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var ( // CLI options, that provide some of the ping functionality
	Count    = flag.Int("c", 3, "number of pings: <= 0 means forever")
	Interval = flag.Duration("i", time.Second, "interval between pings")
	Timeout  = flag.Duration("W", 5*time.Second, "time to wait for a reply")
)

// >> CLI flags help function
func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] host:port\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

// >> should be used in main
func ICMP() {
	// >> reading flags
	flag.Parse()
	if flag.NArg() != 1 { // if host is not given
		fmt.Print("host:port is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// >> getting target host
	target := flag.Arg(0)
	fmt.Println("PING", target)

	// >> @ infinite pings
	if *Count <= 0 { // 0 pings
		fmt.Println("CTRL+C to stop.")
	}

	// >> Pinging target <msg> number of times and
	// logging duration to establish a tcp connection
	msg := 0
	for (*Count <= 0) || (msg < *Count) {
		msg++
		fmt.Print(msg, " ")

		// >> Getting duration from connection attempt to established
		start := time.Now()
		c, err := net.DialTimeout("tcp", target, *Timeout)
		dur := time.Since(start)

		if err != nil {
			fmt.Printf("fail in %s: %v\n", dur, err)
			if nErr, ok := err.(net.Error); !ok || !nErr.Temporary() {
				os.Exit(1)
			}
		} else {
			_ = c.Close()
			fmt.Println(dur)
		}
		time.Sleep(*Interval)
	}

}

// >> ICMP but
// if we want to use ping in our program
func PingTarget(address string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()
	_, err := net.DialTimeout("tcp", address, *Timeout)
	if err != nil {
		return 0, err
	}
	dur := time.Since(start)
	return dur, nil
}
