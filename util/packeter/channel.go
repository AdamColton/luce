package packeter

import (
	"github.com/adamcolton/luce/ds/channel"
)

// UnpackPipe reads chunks of bytes from chs.Rcv, unpacks them with u and sends
// each message on chs.Snd. It returns when chs.Rcv is closed, and closes chs.Snd
// if cls is true.
func UnpackPipe(u Unpacker, chs channel.Pipe[[]byte], cls bool) {
	for msg := range chs.Rcv {
		channel.Slice(u.Unpack(msg), chs.Snd)
	}
	if cls {
		close(chs.Snd)
	}
}

// PackPipe reads messages from chs.Rcv, packs them with p and sends each piece
// on chs.Snd. It returns when chs.Rcv is closed, and closes chs.Snd if cls is
// true.
func PackPipe(p Packer, chs channel.Pipe[[]byte], cls bool) {
	for msg := range chs.Rcv {
		channel.Slice(p.Pack(msg), chs.Snd)
	}
	if cls {
		close(chs.Snd)
	}
}

// Run connects p to a stream that is carried by chs: bytes to write go on
// chs.Snd and bytes that were read arrive on chs.Rcv. It returns a Pipe of
// messages. A message sent on the returned Snd is packed and written to chs.Snd,
// and the messages unpacked from chs.Rcv arrive on the returned Rcv. The
// returned Snd is only set if p is a Packer and chs.Snd is not nil, and the
// returned Rcv is only set if p is an Unpacker and chs.Rcv is not nil. Closing
// the returned Snd closes chs.Snd, and chs.Rcv closing closes the returned Rcv.
func Run(p any, chs channel.Pipe[[]byte]) (out channel.Pipe[[]byte]) {
	if packer, ok := p.(Packer); ok && chs.Snd != nil {
		pp, snd, _ := channel.NewPipe(nil, chs.Snd)
		out.Snd = snd
		go PackPipe(packer, pp, true)
	}
	if unpacker, ok := p.(Unpacker); ok && chs.Rcv != nil {
		up, _, rcv := channel.NewPipe(chs.Rcv, nil)
		out.Rcv = rcv
		go UnpackPipe(unpacker, up, true)
	}

	return
}
