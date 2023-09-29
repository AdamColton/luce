package serialbus

// == projects.Code.luce.serialbus ==

// [ ] move serialbus.String
// 	String(in <-chan []byte) <-chan string

// String converts []byte to string on a channel. The channel it returns is closed
// after in is closed and everything on it has been converted.
func String(in <-chan []byte) <-chan string {
	out := make(chan string, len(in))
	go func() {
		for b := range in {
			out <- string(b)
		}
		close(out)
	}()
	return out
}
