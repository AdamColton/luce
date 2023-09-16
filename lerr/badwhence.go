package lerr

import "strconv"

// ErrBadWhence indicates that a whence value other than SeekStart(0),
// SeekCurrent(1) or SeekEnd(2) was used.
type ErrBadWhence int

// Whence returns ErrBadWhence if i is a value other than SeekStart(0),
// SeekCurrent(1) or SeekEnd(2).
func Whence(i int) error {
	if i < 0 || i > 2 {
		return ErrBadWhence(i)
	}
	return nil
}

// Error fulfills the error type and indicates what bad value was used.
func (e ErrBadWhence) Error() string {
	return "lerr: Seek whence value should be io.SeekStart (0), io.SeekCurrent (1) or io.SeekEnd (2), got:" + strconv.Itoa(int(e))
}
