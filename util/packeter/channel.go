package packeter

import "github.com/adamcolton/luce/ds/channel"

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
