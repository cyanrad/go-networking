package data

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestHttpListener(t *testing.T) {

	// >> server init
	listener, err := net.Listen("tcp", "192.168.0.189:8080")
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
		buf := make([]byte, 1<<19) // 512 KB

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
	}()

	time.Sleep(time.Second * 10)
	listener.Close()
}
