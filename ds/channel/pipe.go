package channel

// Pipe holds a pair of channels. This can be a useful structure for embedding.
// It is also useful for creating routines that operate on values from Rcv and
// put the results on Snd.
type Pipe[T any] struct {
	Rcv <-chan T
	Snd chan<- T
}

// NewPipe creates a Pipe. If rcv is nil, a channel is created for the Pipe's Rcv
// and its sending end is returned as rcvPair. If snd is nil, a channel is
// created for the Pipe's Snd and its receiving end is returned as sndPair. When
// a channel is provided, its pair is nil.
func NewPipe[T any](rcv <-chan T, snd chan<- T) (pipe Pipe[T], rcvPair chan<- T, sndPair <-chan T) {
	if rcv == nil {
		ch := make(chan T)
		rcv = ch
		rcvPair = ch
	}
	if snd == nil {
		ch := make(chan T)
		snd = ch
		sndPair = ch
	}
	pipe = Pipe[T]{
		Snd: snd,
		Rcv: rcv,
	}
	return
}
