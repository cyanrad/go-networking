package data

import (
	"bufio"   // connection scanner
	"bytes"   // used for byte manip at comma_splitter
	"net"     // listeners & dialer
	"reflect" //used for testing the correctness of the data we recieved
	"testing"
)

const payload = "The bigger, the interface, the weaker the abstraction."

// >> function that deletes the '\r' character from the bytes
func dropCarriageReturn(data []byte) []byte {
	// if data.len > 0 && last character is '\r'
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1] // cutting the \r character
	} //otherwise do nothing
	return data
}

// >> Splits the byte slice at ,
// note: removes the '\r' character
// used for scanner.Split(comma_splitter)
func comma_splitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, ','); i >= 0 {
		// +2 advancing so that we skip the space character
		// +1 if we don't want to skip that
		return i + 2, dropCarriageReturn(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCarriageReturn(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func TestScanner(t *testing.T) {
	// >> generic listener
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
		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	// >> generic dialer
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// >> the scanner
	scanner := bufio.NewScanner(conn) // getting a connection scanner, to read (interface of io.Read)
	scanner.Split(comma_splitter)     // splitting the splitter function (split at ", ")
	var words []string
	for scanner.Scan() { // reading the data until we reach a delimeter character (", ")
		// appending the string we read to the data array
		words = append(words, scanner.Text())
	}

	// >> scanner err handling
	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}
	t.Log(words)

	// >> checking if the message we recived is the one we expect
	expected := []string{"The bigger", "the interface", "the weaker the abstraction."}
	t.Log(expected)
	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}
	t.Logf("Scanned words: %#v", words)
}
