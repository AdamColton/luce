package reflector_test

import (
	"fmt"

	"github.com/adamcolton/luce/util/reflector"
)

func ExampleType() {
	t := reflector.Type[string]()
	fmt.Println("t is reflect.Type on", t.String())
	// Output: t is reflect.Type on string
}
