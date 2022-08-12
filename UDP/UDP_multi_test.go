package udp

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
)

func TestListenPacketUDP(t *testing.T) {
	// >> Creating echo server
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	// >> creating node
	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	// >> creating node
	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// >> writing a message to client node
	interrupt := []byte("pardon me")
	n, err := interloper.WriteTo(interrupt, client.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	_ = interloper.Close()

	// >> Checking for sent byte count
	if l := len(interrupt); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	// >> writing from client to the echo server
	ping := []byte("ping")
	_, err = client.WriteTo(ping, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)
	// >> client node reading
	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))

	// >> if recieved byte count is not the same as sent
	if !bytes.Equal(interrupt, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", interrupt, buf[:n])
	}
	// >> if the expected sender is incorrect
	if addr.String() != interloper.LocalAddr().String() {
		t.Errorf("expected message from %q; actual sender is %q",
			interloper.LocalAddr(), addr)
	}

	// >> interloper node reading
	// since the echo server echos the ping to all nodes, it gets echoed to
	// the interloper
	n, addr, err = client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))

	if !bytes.Equal(ping, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", ping, buf[:n])
	}
	if addr.String() != serverAddr.String() {
		t.Errorf("expected message from %q; actual sender is %q",
			serverAddr, addr)
	}
}
