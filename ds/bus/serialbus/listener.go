package serialbus

// TODO: move this...somewhere
// String converts []byte to string on a channel.
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
