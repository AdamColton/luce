package packeter

import "github.com/adamcolton/luce/ds/channel"

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
