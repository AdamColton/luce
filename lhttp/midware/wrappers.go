package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/util/linject"
)

type Initilizer interface {
	Initilize(reflect.Type) DataInserter
}

type DataInserter interface {
	Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) (func(), error)
}

type wrappedInitilizer struct {
	Initilizer
}

func (wi wrappedInitilizer) Initilize(fn linject.Func) linject.DataInserter {
	di := wi.Initilizer.Initilize(fn.DataType())
	return wrappedDataInserter{di}
}

type wrappedDataInserter struct {
	DataInserter
}

func (wdi wrappedDataInserter) Insert(args []reflect.Value) (callback func(), err error) {
	w := args[0].Interface().(http.ResponseWriter)
	r := args[1].Interface().(*http.Request)
	d := args[2]
	return wdi.DataInserter.Insert(w, r, d)
}
