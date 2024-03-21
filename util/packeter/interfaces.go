package packeter

type Packer interface {
	Pack([]byte) [][]byte
}

type Unpacker interface {
	Unpack([]byte) [][]byte
}
