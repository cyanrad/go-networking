package udp

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
)

func TestDialUDP(t *testing.T) {
	// >> creating echo server
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	// >> creating client using dial
	// getting the UDP net.Conn obj
	// note that no handshake is done, so the server recieves no data.
	client, err := net.Dial("udp", serverAddr.String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	// >> Creating the interloper to interrupt.
	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// >> Sending the Interrupt message
	interrupt := []byte("pardon me")
	n, err := interloper.WriteTo(interrupt, client.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}

	// >> clsing interloper
	_ = interloper.Close()
	if l := len(interrupt); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	// >> sending the ping message to the server
	ping := []byte("ping")
	_, err = client.Write(ping)
	if err != nil {
		t.Fatal(err)
	}

	// >> reading data
	buf := make([]byte, 1024)
	n, err = client.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	// >> Error handeling
	if !bytes.Equal(ping, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", ping, buf[:n])
	}
	err = client.SetDeadline(time.Now().Add(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	// >> if unexpected packet arrives
	_, err = client.Read(buf)
	if err == nil {
		t.Fatal("unexpected packet")
	}
}
