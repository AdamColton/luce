package lerr

import (
	"fmt"
	"strings"
)

// DefaultCtxSeperator is used to separate the Ctx description from the inner
// error when Ctx.Seperator is empty.
var DefaultCtxSeperator = ": "

// Ctx adds context to an error.
type Ctx struct {
	Innr      error
	Desc      string
	Seperator string
}

// Error fulfills error. It joins Desc and the inner error with Seperator, or
// with DefaultCtxSeperator if Seperator is empty.
func (err Ctx) Error() string {
	sep := err.Seperator
	if sep == "" {
		sep = DefaultCtxSeperator
	}
	return strings.Join([]string{err.Desc, err.Innr.Error()}, sep)
}

// Wrap an error to provide context. If the error is nil, nil will be returned.
func Wrap(err error, format string, data ...interface{}) error {
	if err == nil {
		return nil
	}
	return Ctx{
		Innr: err,
		Desc: fmt.Sprintf(format, data...),
	}
}
