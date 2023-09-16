package lerr

import "strconv"

type ErrBadWhence int

func (e ErrBadWhence) Error() string {
	return "lerr: Seek whence value should be io.SeekStart (0), io.SeekCurrent (1) or io.SeekEnd (1), got:" + strconv.Itoa(int(e))
}
