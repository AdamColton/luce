package midware

import (
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/adamcolton/luce/util/reflector/ltype/httptype"
)

// HttpHandlerType matches a func type with three arguments whose first is an
// http.ResponseWriter and second is a *http.Request — the shape required of
// a MidwareFunc.
var (
	HttpHandlerType = filter.NumInEq(3).
		And(filter.InType(0, httptype.ResponseWriter)).
		And(filter.InType(1, httptype.Request))
)

// NewField wraps fsi as a linject.Field for fieldName, constrained to
// HttpHandlerType so it can only be used to populate a MidwareFunc's data.
func NewField(fsi linject.FieldInitilizer, fieldName string) linject.Field {
	fi := linject.NewField(fsi, fieldName)
	fi.FuncType = HttpHandlerType
	return fi
}
