package lerr

import (
	"fmt"
	"strings"
)

// SliceErrRecord represents an error comparing two slices. An index less than
// 0 is treated as an error on the underlying comparison and will be printed
// without the index. For instance, if the lengths of the slices are not the
// same.
type SliceErrRecord struct {
	// Index of the reference slice
	Index int
	Err   error
}

// SliceErrs collects SliceErrRecord and treats them as a single error.
type SliceErrs []SliceErrRecord

// MaxSliceErrs limits the number of errors that will be reported by
// SliceErrs.Error
var MaxSliceErrs = 10

// Error fulfills the error interface.
func (e SliceErrs) Error() string {
	var out []string
	ln := len(e)
	if ln > MaxSliceErrs {
		out = make([]string, MaxSliceErrs+1)
		out[MaxSliceErrs] = fmt.Sprintf("Omitting %d more", ln-MaxSliceErrs)
		ln = MaxSliceErrs
	} else {
		out = make([]string, ln)
	}

	for i := 0; i < ln; i++ {
		r := e[i]
		if r.Index < 0 {
			out[i] = r.Err.Error()
		} else {
			out[i] = fmt.Sprintf("\t%d: %s", r.Index, r.Err.Error())
		}
	}

	return strings.Join(out, "\n")
}

// Append a SliceErrRecord to SliceErrs. If err is nil it will not be appended.
func (e SliceErrs) Append(idx int, err error) SliceErrs {
	if err == nil {
		return e
	}
	return append(e, SliceErrRecord{
		Index: idx,
		Err:   err,
	})
}

// AppendF uses fmt.Errorf to append a SliceErrRecord to SliceErrs .
func (e SliceErrs) AppendF(idx int, format string, args ...interface{}) SliceErrs {
	return append(e, SliceErrRecord{
		Index: idx,
		Err:   fmt.Errorf(format, args...),
	})
}

func (e SliceErrs) Cast() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

// NewSliceErrs creates an instance of SliceErrs by calling the provided func
// for every value up to Min(lnExpected, lnActual). If lnExpected and lnActual
// are not equal and instance of LenMismatch will be added to the start of the
// SliceErrs with an index of -1. If lnActual == -1 the length check is ignored.
func NewSliceErrs(lnExpected, lnActual int, fn func(int) error) error {
	var out SliceErrs
	var ln int
	if lnActual >= 0 {
		var lnErr error
		ln, _, lnErr = NewLenMismatch(lnExpected, lnActual)
		out = out.Append(-1, lnErr)
	} else {
		ln = lnExpected
	}
	for i := 0; i < ln; i++ {
		out = out.Append(i, fn(i))
	}
	return out.Cast() // YEEEEEEEAHHHHH
}
