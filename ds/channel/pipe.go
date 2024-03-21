package channel

// Pipe holds a pair of channels. It is intended for embedding.
type Pipe[T any] struct {
	Rcv <-chan T
	Snd chan<- T
}

func NewPipe[T any](rcv <-chan T, snd chan<- T) (pipe Pipe[T], retSnd chan<- T, retRcv <-chan T) {
	if rcv == nil {
		ch := make(chan T)
		rcv = ch
		retSnd = ch
	}
	if snd == nil {
		ch := make(chan T)
		snd = ch
		retRcv = ch
	}
	pipe = Pipe[T]{
		Snd: snd,
		Rcv: rcv,
	}
	return
}
