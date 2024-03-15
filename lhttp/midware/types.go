package midware

import (
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector/ltype/httptype"
)

var (
	HttpHandlerType = filter.NumInEq(3).
		And(filter.InType(0, httptype.ResponseWriter)).
		And(filter.InType(1, httptype.Request))
)
