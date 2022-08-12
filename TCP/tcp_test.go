package TCP

import (
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

// >> Creating a tcp server, and checking if it binds correctly
func TestListener(t *testing.T) {
	// >> ip string to compare
	// port is random, but using 11211 instead of 0
	// to have defined port
	ipString := "127.0.0.1:11211"

	// creating a net.listener (aka tcp ipv4 listener on localhost)
	t.Log("creating server")
	listener, err := net.Listen("tcp4", ipString) // :0 = random port
	require.NoError(t, err)
	defer listener.Close() // Closing the server when done

	// logs address of server,
	// when the server binding to an address
	require.Equal(t, ipString, listener.Addr().String())
	t.Log("bound to", listener.Addr())
}

// >> tests basic client-server connection,
// and transmission of data
func TestDial(t *testing.T) {
	testData := []byte("this is a byte slice")

	// >> Create server on a random port
	t.Log("creating server")
	listener, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	// >> done channel to wait for clients on main thread
	done := make(chan struct{})

	// >> Handling tcp clients
	go func() {
		// >> To continue main thread when done with routine
		// closing the server
		defer func() { done <- struct{}{} }()

		// >> accepting connections to sever
		// and getting their obejct
		t.Log("accepting client connection")
		conn, err := listener.Accept()
		require.NoError(t, err)

		// >> close connections (sends FIN pkt to conn)
		defer func() {
			t.Log("closing connection")
			conn.Close()
			done <- struct{}{}
		}()

		// >> to make sure connection is established first
		done <- struct{}{}

		// >> reads & logs 1024 byte, from socket
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf) // reading
			t.Log("reading")
			// will return io.EOF when client closes
			if err != nil {
				require.EqualError(t, err, io.EOF.Error())
				t.Log("EOF")
				return
			}

			require.Equal(t, buf[:n], testData)
			t.Logf("received: %q", buf[:n])
		}
	}()
	// >> creating tcp clients
	// note: service name can be used (E.g. http)
	t.Log("creating client")
	client, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)

	// >> Closing client (sending FIN packet to the c.Read method)
	<-done
	t.Log("sending data")
	client.Write([]byte(testData))
	t.Log("closing connection from client side")
	client.Close()
	<-done

	// >> Done
	t.Log("closing server")
	listener.Close()
	<-done
}
