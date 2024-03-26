package channel

type Writer struct {
	Ch chan<- []byte
}

func (w Writer) Write(p []byte) (n int, err error) {
	w.Ch <- p
	return len(p), nil
}

// TODO Writer.ReadFrom, Reader.Read, Reader.WriteTo
