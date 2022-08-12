package TCP

import (
	"io"
	"net"
)

func proxyConn(source, destination string) error {
	connSource, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}
	defer connSource.Close()
	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer connDestination.Close()
	// connDestination replies to connSource
	go func() { _, _ = io.Copy(connSource, connDestination) }()
	// connSource messages to connDestination
	_, err = io.Copy(connDestination, connSource)
	return err
}

func proxy(from io.Reader, to io.Writer) error {
	// connecting reader and writer
	fromWriter, fromIsWriter := from.(io.Writer)
	toReader, toIsReader := to.(io.Reader)

	if toIsReader && fromIsWriter {
		// Send replies since "from" and "to" implement the
		// necessary interfaces.
		go func() { _, _ = io.Copy(fromWriter, toReader) }()
	}
	_, err := io.Copy(to, from)
	return err
}
