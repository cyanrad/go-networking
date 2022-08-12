package TCP

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPingerAdvanceDeadline(t *testing.T) {
	done := make(chan struct{})

	// >> Creating server
	t.Log("creating server")
	listener, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	// >> pinger
	// deadline: 			5 seconds
	// ping interval: 		every second
	// from client side: 	recieve 4 pings then an io.EOF
	begin := time.Now()
	go func() {
		// closing the done channel when process is complete
		defer func() { close(done) }()

		// >> accepting conneciton
		t.Log("accepting connection")
		conn, err := listener.Accept()
		require.NoError(t, err)

		// >> pinger control context
		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
			t.Log("closing connection")
			conn.Close()
		}()

		// >> buffered timer reset channel
		resetTimer := make(chan time.Duration, 1)
		resetTimer <- time.Second

		// >> Running pinger
		go pinger(ctx, conn, resetTimer)

		// >> setting connection deadline
		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		require.NoError(t, err)

		buf := make([]byte, 1024)
		for {
			// >> writing to buffer client data
			n, err := conn.Read(buf) // blocks until data is recieved
			if err != nil {
				return
			}
			t.Logf("[%s] %s", //logging read data
				time.Since(begin).Truncate(time.Second), buf[:n])

			// resetting pinger
			resetTimer <- 0

			// resetting deadline
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			require.NoError(t, err)
		}
	}()

	// >> client init
	conn, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	// >> Reading pings
	buf := make([]byte, 1024)
	for i := 0; i < 4; i++ { // read up to four pings
		n, err := conn.Read(buf)
		require.NoError(t, err)
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	// writing to server
	_, err = conn.Write([]byte("PONG!!!")) // should reset the ping timer
	require.NoError(t, err)

	for i := 0; i < 4; i++ { // read up to four more pings
		n, err := conn.Read(buf)
		if err != nil {
			require.EqualError(t, err, io.EOF.Error())
			break
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	// @ finish (checking if time is correct)
	<-done
	end := time.Since(begin).Truncate(time.Second)
	t.Logf("[%s] done", end)
	require.Equal(t, end, 9*time.Second)
}
