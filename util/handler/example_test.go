package handler_test

import (
	"fmt"

	"github.com/adamcolton/luce/util/handler"
)

func ExampleHandler() {
	h, err := handler.New(func(name string) (string, error) {
		if name == "" {
			return "", fmt.Errorf("a name is required")
		}
		return "hello " + name, nil
	})
	fmt.Println(err)

	// The type of the argument tells a caller what to pass to Handle.
	fmt.Println(h.Type())

	out, err := h.Handle("Ada")
	fmt.Println(out, err)

	out, err = h.Handle("")
	fmt.Printf("%q %v\n", out, err)
	// Output:
	// <nil>
	// string
	// hello Ada <nil>
	// "" a name is required
}
