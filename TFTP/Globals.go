package tftp

const (
	DatagramSize = 516              // the maximum supported datagram size
	BlockSize    = DatagramSize - 4 // the DatagramSize minus a 4-byte header
)

type OpCode uint16
type ErrCode uint16

// >> operation codes
const (
	OpRRQ  OpCode = iota + 1
	_             // no WRQ support (write request)
	OpData        // get data
	OpAck         // Ack
	OpErr         // Error
	OpGet         // Exsisting files
)

// >> Error codes
const (
	ErrUnknown ErrCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)
