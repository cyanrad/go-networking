package tftp

import (
	"bytes"
	"errors"
	"log"
	"net"
	"time"
)

type Server struct {
	Payload []byte        // the payload served for all read requests
	Retries uint8         // the number of times to retry a failed transmission
	Timeout time.Duration // the duration to wait for an acknowledgment
}

func (s Server) ListenAndServe(addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	log.Printf("Listening on %s ...\n", conn.LocalAddr())
	return s.Serve(conn)
}

// >> accepts the RRq and starts the process of sending data pkts
// $$$ Feature: the payload should depend on the file requested.
func (s *Server) Serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("nil connection")
	}
	if s.Payload == nil {
		return errors.New("payload is required")
	}
	if s.Retries == 0 {
		s.Retries = 10
	}
	if s.Timeout == 0 {
		s.Timeout = 6 * time.Second
	}
	var rrq ReadReq
	for {
		buf := make([]byte, DatagramSize)
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}
		err = rrq.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("[%s] bad request: %v", addr, err)
			continue
		}
		// @ early return -> the entire payload is not sent
		//
		go s.handle(addr.String(), rrq)
	}
}

// handle sending the data
func (s Server) handle(clientAddr string, rrq ReadReq) {
	log.Printf("[%s] requested file: %s", clientAddr, rrq.Filename)

	// >> creating a client
	// So that the server is not busy with communicating
	// and so we can only recieve data from the client
	conn, err := net.Dial("udp", clientAddr)
	if err != nil {
		log.Printf("[%s] dial: %v", clientAddr, err)
		return
	}
	defer func() { _ = conn.Close() }()

	var ( // >> creating some variables
		ackPkt  Ack
		errPkt  Err
		dataPkt = Data{Payload: bytes.NewReader(s.Payload)}
		buf     = make([]byte, DatagramSize)
	)

	// creating the data packet response
	// the connection ends at n != DatagramSize
	// since when that is the case we have reached the final packet
NEXTPACKET:
	// n declereation will only run once i think
	// n gets updated with the size of the previous dataRq sent
	for n := DatagramSize; n == DatagramSize; {
		data, err := dataPkt.MarshalBinary()
		if err != nil {
			log.Printf("[%s] preparing data packet: %v", clientAddr, err)
			return
		}
		// trying to send the packet
	RETRY:
		for i := s.Retries; i > 0; i-- {
			n, err = conn.Write(data) // send the data packet
			if err != nil {
				log.Printf("[%s] write: %v", clientAddr, err)
				return
			}

			// wait for the client's ACK packet
			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))
			_, err = conn.Read(buf) // reading ACK
			if err != nil {
				// if err is timeout
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					log.Println("Testing")
					continue RETRY
				}
				log.Printf("[%s] waiting for ACK: %v", clientAddr, err)
				return
			}

			switch {
			case ackPkt.UnmarshalBinary(buf) == nil: // > correct ack
				if uint16(ackPkt) == dataPkt.Block {
					// received ACK; send next data packet
					continue NEXTPACKET
				}
			case errPkt.UnmarshalBinary(buf) == nil: // > err
				log.Printf("[%s] received error: %v",
					clientAddr, errPkt.Message)
				return
			default: // > unknown
				log.Printf("[%s] bad packet", clientAddr)
				return
			}
		}
		// >> at retry count too high
		log.Printf("[%s] exhausted retries", clientAddr)
		return
	}
	log.Printf("[%s] sent %d blocks", clientAddr, dataPkt.Block)
}
