package serial_test

import (
	"encoding/json"
	"io"

	"github.com/adamcolton/luce/lerr"
)

type person struct {
	Name string
	Age  int
}

func mockSerialize(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

func errSerializeFn(i interface{}, w io.Writer) error {
	return errSerialize
}

func mockDeserialize(i interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

var (
	jsonStr    = "{\"Name\":\"Adam\",\"Age\":39}\n"
	testPerson = person{
		Name: "Adam",
		Age:  39,
	}
	errSerialize = lerr.Str("serialize error")
)
