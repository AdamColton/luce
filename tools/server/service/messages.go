package service

// SocketOpened notifies a service that a client socket, identified by ID,
// has opened.
type SocketOpened struct {
	ID uint32
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketOpened) TypeID32() uint32 {
	return 1046109042
}

// SocketClose notifies a service that the client socket identified by ID
// has closed.
type SocketClose struct {
	ID uint32
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketClose) TypeID32() uint32 {
	return 3196974518
}

// SocketMessage carries a message Body from the client socket identified by
// ID.
type SocketMessage struct {
	ID   uint32
	Body []byte
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketMessage) TypeID32() uint32 {
	return 3196974518
}
