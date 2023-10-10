// Package procbus holds inline processors for chan []byte busses
package procbus

import "strings"

// Delim reads from a byte channel, treating the data as a string and queues up
// data until the specified delimitor is reached. Then the entire queue is
// written as a single message.
func Delim(in <-chan []byte, delim rune) <-chan []byte {
	out := make(chan []byte)
	go runDelim(in, out, delim)
	return out
}

// TODO: circular buffer
func runDelim(in <-chan []byte, out chan<- []byte, delim rune) {
	var msg []byte
	for b := range in {
		str := string(b)
		idx := strings.IndexRune(str, delim)

		if idx < 0 {
			msg = append(msg, b...)
		} else {
			msg = append(msg, str[:idx]...)
			out <- msg

			rest := []byte(str[idx+1:])
			msg = make([]byte, len(rest))
			copy(msg, rest)
		}
	}

	if len(msg) > 0 {
		out <- msg
	}
}
