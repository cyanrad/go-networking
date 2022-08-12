package data

import (
	"net"
	"reflect"
	"testing"
)

func TestPayloads(t *testing.T) {
	// >> payloads
	b1 := Binary("Clear is better than clever.")
	b2 := Binary("Don't panic.")
	s1 := String("Errors are values.")
	payloads := []Payload{&b1, &s1, &b2}

	// >> generic Listener
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		// Writing
		for _, p := range payloads {
			_, err = p.WriteTo(conn)
			if err != nil {
				t.Error(err)
				break
			}
		}
	}()

	// >> generic Dialer
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < len(payloads); i++ {
		actual, err := decode(conn) // getting the actual data
		if err != nil {
			t.Fatal(err)
		}
		if expected := payloads[i]; !reflect.DeepEqual(expected, actual) {
			t.Errorf("value mismatch: %v != %v", expected, actual)
			continue
		}
		t.Logf("[%T] %[1]q", actual)
	}
}
