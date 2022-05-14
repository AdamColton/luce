package lhttp

// Socket is intended to be a websocket
type Socket struct {
	MessageReaderWriter
}

// NewSocket creates a Socket. It is intended to be used with a websocket.
func NewSocket(socket MessageReaderWriter) Socket {
	return Socket{
		MessageReaderWriter: socket,
	}
}

// RunReader reads from the socket until it encounters an error and writes each
// msg to the "from" channel.
func (socket Socket) RunReader(from chan<- []byte) {
	for {
		_, msg, err := socket.ReadMessage()
		if err != nil {
			break
		}
		from <- msg
	}
	close(from)
}

// RunSender pulls messages off the "to" channel and writes them to the
// websocket.
func (socket Socket) RunSender(to <-chan []byte) {
	for msg := range to {
		err := socket.WriteMessage(1, msg)
		if err != nil {
			break
		}
	}
}
