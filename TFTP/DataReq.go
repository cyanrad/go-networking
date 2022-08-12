package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Data struct {
	Block   uint16
	Payload io.Reader // the payload can be large
}

// >> Data response
// create data request
// returns 516 bytes, will syphon these bytes from the payload
func (d *Data) MarshalBinary() ([]byte, error) {
	// >> creatin byte slice
	b := new(bytes.Buffer)
	b.Grow(DatagramSize)

	d.Block++                                        // block numbers increment from 1
	err := binary.Write(b, binary.BigEndian, OpData) // write operation code
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, d.Block) // write block number
	if err != nil {
		return nil, err
	}

	// payload
	// write up to BlockSize worth of bytes
	_, err = io.CopyN(b, d.Payload, BlockSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *Data) UnmarshalBinary(p []byte) error {
	// >> Sanity Check
	if l := len(p); l < 4 || l > DatagramSize {
		return errors.New("invalid DATA")
	}

	// >> getting and checking opcode
	var opcode OpCode
	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &opcode)
	if err != nil || opcode != OpData {
		return errors.New("invalid DATA")
	}

	// >> Getting block count
	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block)
	if err != nil {
		return errors.New("invalid DATA")
	}

	// getting payload
	d.Payload = bytes.NewBuffer(p[4:])
	return nil
}
