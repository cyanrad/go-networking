package UDS

import (
	"context"
	"net"
)

// >> creating an echo server
// you can pass any stream oriented network type (TCP, UDS, unixpacket)
// ctx, network type, network address
func streamingEchoServer(ctx context.Context, network string,
	addr string) (net.Addr, error) {
	// >> Creating server
	s, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	// >> running the echo server
	go func() {

		// >> Closing server at ctx done
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		for { // loop so that it can accept more than one connection
			// >> Accepting connections
			conn, err := s.Accept()
			if err != nil {
				return
			}

			go func() { // after running this the func gets detached
				defer func() { _ = conn.Close() }()
				for {
					// >> reading data
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					// >> Echoing it
					_, err = conn.Write(buf[:n])
					if err != nil {
						return
					}
				}
			}() // single client go finish
		} // Acceptor loop finish
	}()

	return s.Addr(), nil
}
