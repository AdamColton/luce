package channel

// == projects.Code.luce.channel ==

// Writer uses a []byte channel to fulfill io.Writer
type Writer struct {
	Ch chan<- []byte
}

// Write fulfills io.Writer, sending data on the channel and blocking until it is
// received. It sends data itself, not a copy, so the receiver shares its memory
// with the caller. That does not follow the io.Writer rule that Write must not
// retain the slice it is given.
func (w Writer) Write(data []byte) (n int, err error) {
	w.Ch <- data
	return len(data), nil
}

// [ ] channel.Writer.ReadFrom
// [ ] channel.Reader.Read
// [ ] channel.Reader.WriteTo
