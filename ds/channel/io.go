package channel

import "slices"

// == projects.Code.luce.channel ==
// [ ] reduce duplication
// some of this logic overlaps with ds/bus/iobus
// is it possible to reduce the duplication?

// Writer uses a []byte channel to fulfill io.Writer
type Writer struct {
	Ch chan<- []byte
}

// Write fulfills io.Writer, sending a copy of data on the channel and blocking
// until it is received. The copy means the caller is free to reuse data as soon
// as Write returns, as the io.Writer contract requires.
func (w Writer) Write(data []byte) (n int, err error) {
	w.Ch <- slices.Clone(data)
	return len(data), nil
}

// [ ] channel.Writer.ReadFrom
// [ ] channel.Reader.Read
// [ ] channel.Reader.WriteTo
