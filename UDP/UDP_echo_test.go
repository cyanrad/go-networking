package udp

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"testing"
)

// net.Addr is used for addressing messages to the echo server.
func echoServerUDP(ctx context.Context, addr string) (net.Addr, error) {
	// Listening to packets on the specified address.
	// s type is net.PacketConn
	// very much the same as net.Listen, but for UDP
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("binding to udp %s: %w", addr, err)
	}

	// managing the echoing of packets.
	// once the caller cancles the context the server closes.
	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		buf := make([]byte, 1024)
		for {
			// reading
			// you rely on the address to know which node you are comm with.
			n, clientAddr, err := s.ReadFrom(buf) // client to server
			if err != nil {
				return
			}

			// Echo any packets sent, returns the number of bytes written.
			_, err = s.WriteTo(buf[:n], clientAddr) // server to client
			if err != nil {
				return
			}
		}
	}()
	return s.LocalAddr(), nil
}

func TestEchoServerUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()
	msg := []byte("Eat my ass")
	_, err = client.WriteTo(msg, serverAddr)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
	if addr.String() != serverAddr.String() {
		t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
	}
	if !bytes.Equal(msg, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", msg, buf[:n])
	}
}
