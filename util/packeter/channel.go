package packeter

import (
	"github.com/adamcolton/luce/ds/channel"
)

func UnpackPipe(u Unpacker, chs channel.Pipe[[]byte], cls bool) {
	for msg := range chs.Rcv {
		channel.Slice(u.Unpack(msg), chs.Snd)
	}
	if cls {
		close(chs.Snd)
	}
}

func PackPipe(p Packer, chs channel.Pipe[[]byte], cls bool) {
	for msg := range chs.Rcv {
		channel.Slice(p.Pack(msg), chs.Snd)
	}
	if cls {
		close(chs.Snd)
	}
}

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
