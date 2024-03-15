package midware

import (
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/adamcolton/luce/util/reflector/ltype/httptype"
)

var (
	HttpHandlerType = filter.NumInEq(3).
		And(filter.InType(0, httptype.ResponseWriter)).
		And(filter.InType(1, httptype.Request))
)

func NewField(fsi linject.FieldInitilizer, fieldName string) linject.Field {
	fi := linject.NewField(fsi, fieldName)
	fi.FuncType = HttpHandlerType
	return fi
}
