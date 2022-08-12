package TCP

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDeadline(t *testing.T) {
	bytes := []byte("test")
	duration := time.Second * 5
	sync := make(chan struct{}) // used to control read flow

	// >> Creating server
	t.Log("creating server")
	listener, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	// >> Accepting connections
	go func() {
		t.Log("connecting")
		conn, err := listener.Accept()
		require.NoError(t, err)

		defer func() {
			t.Log("closing connection")
			conn.Close()
			close(sync)
		}()

		// >> Setting deadline
		err = conn.SetDeadline(time.Now().Add(duration))
		require.NoError(t, err)

		// >> Reading data
		// should fails -> client's write function is stopped by the sync channel
		// future reads will result in a timeout error
		buf := make([]byte, len(bytes))
		// block until Data is recived
		t.Log("server reading (should timeout)")
		_, err = conn.Read(buf) //should return an error after 5 sec

		nErr, ok := err.(net.Error) // >> Verifying that the error is a timeout
		require.True(t, ok)
		require.True(t, nErr.Timeout())

		// functionality can be restored with pushing the deadline
		sync <- struct{}{}
		err = conn.SetDeadline(time.Now().Add(duration))
		require.NoError(t, err)

		// this read will succseed since we extended the deadline
		_, err = conn.Read(buf)
		t.Log("server reading (works)")
		t.Log(string(buf))
		require.Equal(t, buf, bytes)
		require.NoError(t, err)
	}()

	// >> Creating client
	t.Log("creating client")
	client, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)
	defer client.Close()

	<-sync // >> to force a timeout from server

	// >> writing to server
	t.Log("writing to server")
	_, err = client.Write(bytes)
	require.NoError(t, err)

	// >> waiting for server to close
	buf := make([]byte, 1)
	//sends data @ closing connection
	data, _ := client.Read(buf) // blocked until remote node sends data
	t.Log("reading data from client:", data)
}
