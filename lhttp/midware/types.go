package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/adamcolton/luce/util/reflector/ltype/httptype"
)

var (
	HttpHandlerType = filter.NumInEq(3).
		And(filter.InType(0, httptype.ResponseWriter)).
		And(filter.InType(1, httptype.Request))
)

func GetWR(args []reflect.Value) (w http.ResponseWriter, r *http.Request) {
	w = args[0].Interface().(http.ResponseWriter)
	r = args[1].Interface().(*http.Request)
	return
}

func NewFieldInitilizer(fsi linject.FieldSetterInitilizer, fieldName string) linject.FieldInitilizer {
	fi := linject.NewFieldInitilizer(fsi, fieldName)
	fi.FuncType = HttpHandlerType
	return fi
}
