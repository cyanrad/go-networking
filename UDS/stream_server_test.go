package UDS

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

// >> testing the UDS echo server
func TestEchoServerUnix(t *testing.T) {
	// creating a temp directory
	dir, err := ioutil.TempDir("", "echo_unix")
	if err != nil {
		t.Fatal(err)
	}

	// >> close dir when done
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(rErr)
		}
	}()

	ctx, close := context.WithCancel(context.Background())
	// >> getting path for socket file
	// joinging the dir with the actual file
	// so the file path will be: dirName/process_ID.sock
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))

	// >> creating the echo server
	rAddr, err := streamingEchoServer(ctx, "unix", socket)
	if err != nil {
		t.Fatal(err)
	}

	// >> changes the mode of the file
	// prolly used for linux
	err = os.Chmod(socket, os.ModeSocket|0666)
	if err != nil {
		t.Fatal(err)
	}

	// >> creating client
	// we do that by dialing the server by its file path
	conn, err := net.Dial("unix", rAddr.String())
	if err != nil {
		t.Fatal(err)
	}

	// >> writing ping messages
	defer func() { _ = conn.Close() }()
	msg := []byte("ping")
	for i := 0; i < 3; i++ { // write 3 "ping" messages
		_, err = conn.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	// >> Reading the echoing packets from the server
	buf := make([]byte, 1024)
	n, err := conn.Read(buf) // read once from the server
	if err != nil {
		t.Fatal(err)
	}

	// >> testing for correctness of messages
	expected := bytes.Repeat(msg, 3)
	if !bytes.Equal(expected, buf[:n]) {
		t.Fatalf("expected reply %q; actual reply %q", expected,
			buf[:n])
	}
	close()
}
