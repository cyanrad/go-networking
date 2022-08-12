package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

type ReadReq struct {
	Filename string
	Mode     string
}

// Although not used by our server, a client would make use of this method.
func (q ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if q.Mode != "" {
		mode = q.Mode
	}
	// operation code + filename + 0 byte + mode + 0 byte
	cap := 2 + 2 + len(q.Filename) + 1 + len(q.Mode) + 1
	b := new(bytes.Buffer)
	b.Grow(cap)
	err := binary.Write(b, binary.BigEndian, OpRRQ) // write operation code
	if err != nil {
		return nil, err
	}
	_, err = b.WriteString(q.Filename) // write filename
	if err != nil {
		return nil, err
	}
	err = b.WriteByte(0) // write 0 byte
	if err != nil {
		return nil, err
	}
	_, err = b.WriteString(mode) // write mode
	if err != nil {
		return nil, err
	}
	err = b.WriteByte(0) // write 0 byte
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// >> read request binary
func (q *ReadReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p) // >> reading binary

	// >> Reading OpCode
	var code OpCode
	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}
	if code != OpRRQ { // > Checking validity
		return errors.New("invalid RRQ")
	}

	// >> Reading Filename
	q.Filename, err = r.ReadString(0)
	if err != nil {
		return errors.New("invalid RRQ")
	}
	q.Filename = strings.TrimRight(q.Filename, "\x00") // remove the 0-byte
	if len(q.Filename) == 0 {
		return errors.New("invalid RRQ")
	}

	// >> Reading Mode
	q.Mode, err = r.ReadString(0)
	if err != nil {
		return errors.New("invalid RRQ")
	}
	q.Mode = strings.TrimRight(q.Mode, "\x00") // remove the 0-byte
	if len(q.Mode) == 0 {
		return errors.New("invalid RRQ")
	}
	actual := strings.ToLower(q.Mode) // > Ensuring Octet mode is selected
	if actual != "octet" {
		return errors.New("only binary transfers supported")
	}
	return nil
}
