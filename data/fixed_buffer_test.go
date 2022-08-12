package data

import (
	"io"
	"net"
	"testing"
)

func TestReadIntoBuffer(t *testing.T) {
	// >> payload init
	var payload string = "p..penis communication"
	//_, err := rand.Read([]byte(payload)) // generate a random payload
	//if err != nil {
	//	t.Fatal(err)
	//}

	// >> server init
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		// >> accepting connection
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer conn.Close()

		// >> sending payload to connection
		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	// >> Client
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1<<19) // 512 KB

	// >> reading data from server
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
		t.Logf("read %d bytes", n) // buf[:n] is the data read from conn
		t.Log(string(buf))
	}
	conn.Close()
}
