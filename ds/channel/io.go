package channel

// Writer uses a []byte channel to fulfill io.Writer
type Writer struct {
	Ch chan<- []byte
}

// Write fulfills io.Writer, writing the data to the channel.
func (w Writer) Write(data []byte) (n int, err error) {
	w.Ch <- data
	return len(data), nil
}

// TODO Writer.ReadFrom, Reader.Read, Reader.WriteTo
