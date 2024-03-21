package packeter

// Packer takes a message and returns the pieces to write to a stream, so that
// the receiver can tell where the message ends. The pieces are meant to be
// written in order.
type Packer interface {
	Pack([]byte) [][]byte
}

// Unpacker takes the next chunk of bytes read from a stream and returns the
// messages that it completes. A chunk may complete no messages or several, and
// the data of an unfinished message is kept until the next call.
type Unpacker interface {
	Unpack([]byte) [][]byte
}
