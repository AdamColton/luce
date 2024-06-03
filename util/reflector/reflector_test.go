package reflector_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func ExampleType() {
	t := reflector.Type[string]()
	fmt.Println("t is reflect.Type on", t.String())
	// Output: t is reflect.Type on string
}

func TestToType(t *testing.T) {
	s := "string"
	st := reflect.TypeOf(s)

	assert.Equal(t, st, reflector.ToType(s))
	assert.Equal(t, st, reflector.ToType(st))
}

func ExampleToType() {
	t := reflector.ToType("test")
	fmt.Println("t is reflect.Type on", t.String())

	t2 := reflector.ToType(t)
	fmt.Println("t2 is reflect.Type on", t2.String())
	// Output: t is reflect.Type on string
	// t2 is reflect.Type on string
}

func TestToValue(t *testing.T) {
	s := "foo"
	sv := reflect.ValueOf(s)
	str := reflector.ToValue(s).Type().String()
	assert.Equal(t, "string", str)

	assert.Equal(t, sv.Interface(), reflector.ToValue(s).Interface())
	assert.Equal(t, sv, reflector.ToValue(sv))
}

func ExampleToValue() {
	v := reflector.ToValue("test")
	fmt.Println("v is reflect.Value on", v.Kind())

	v2 := reflector.ToValue(v)
	fmt.Println("v2 is reflect.Value on", v2.Kind())
	// Output: v is reflect.Value on string
	// v2 is reflect.Value on string
}

func TestReturnsErrCheck(t *testing.T) {
	tt := map[string]struct {
		fn       any
		args     []reflect.Value
		expected error
	}{
		"no-returns": {
			fn: func() {
			},
			expected: nil,
		},
		"one-return-no-error": {
			fn: func(str string) string {
				return "hello"
			},
			args:     []reflect.Value{reflect.ValueOf("hi")},
			expected: nil,
		},
		"one-return-is-error": {
			fn: func(str string) error {
				return fmt.Errorf("hello")
			},
			args:     []reflect.Value{reflect.ValueOf("hi")},
			expected: fmt.Errorf("hello"),
		},
		"one-return-is-nil-error": {
			fn: func(str string) error {
				return nil
			},
			args:     []reflect.Value{reflect.ValueOf("hi")},
			expected: nil,
		},
		"two-returns-is-error": {
			fn: func(str string) (string, error) {
				return "hello", fmt.Errorf("goodbye")
			},
			args:     []reflect.Value{reflect.ValueOf("hi")},
			expected: fmt.Errorf("goodbye"),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			vfn := reflect.ValueOf(tc.fn)
			err := reflector.ReturnsErrCheck(vfn.Call(tc.args))
			assert.Equal(t, tc.expected, err)
		})
	}
}

func ExampleReturnsErrCheck() {
	v := reflect.ValueOf(func(i int) (int, error) {
		if i > 0 {
			return i + 10, nil
		}
		return 0, fmt.Errorf("i should be > 0, got: %d", i)
	})

	args := []reflect.Value{
		reflect.ValueOf(10),
	}
	got := v.Call(args)
	err := reflector.ReturnsErrCheck(got)
	fmt.Println(err)

	args[0] = reflect.ValueOf(-1)
	got = v.Call(args)
	err = reflector.ReturnsErrCheck(got)
	fmt.Println(err)

	// Output:
	// <nil>
	// i should be > 0, got: -1
}

func TestIsNil(t *testing.T) {
	var strPtr *string

	v := reflect.ValueOf(strPtr)
	assert.True(t, reflector.IsNil(v))

	v = reflect.ValueOf(123)
	assert.False(t, reflector.IsNil(v))
}

func ExampleIsNil() {
	str := "test"
	strPtr := &str
	v := reflect.ValueOf(strPtr)
	fmt.Println(reflector.IsNil(v))

	strPtr = nil
	v = reflect.ValueOf(strPtr)
	fmt.Println(reflector.IsNil(v))

	v = reflect.ValueOf(123)
	// calling v.IsNil() would panic
	fmt.Println(reflector.IsNil(v))

	// Output:
	// false
	// true
	// false
}

func TestMake(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	pt := reflector.Type[*Person]()
	pp := reflector.Make(pt).Interface().(*Person)
	pp.Name = "Adam"

	p := reflector.Make(pt.Elem()).Interface().(Person)
	p.Name = "Lauren"
}

func TestCanElem(t *testing.T) {
	tt := make(map[reflect.Kind]bool, int(reflect.UnsafePointer))
	for k := reflect.Kind(1); k <= reflect.UnsafePointer; k++ {
		tt[k] = false
	}
	tt[reflect.Array] = true
	tt[reflect.Chan] = true
	tt[reflect.Map] = true
	tt[reflect.Pointer] = true
	tt[reflect.Slice] = true
	for n, tc := range tt {
		t.Run(n.String(), func(t *testing.T) {
			assert.Equal(t, tc, reflector.CanElem(n))
		})
	}
}

func TestElem(t *testing.T) {
	tt := map[reflect.Type]bool{
		reflector.Type[string]():         false,
		reflector.Type[[]int]():          true,
		reflector.Type[int]():            false,
		reflector.Type[map[int]string](): true,
		nil:                              false,
	}

	name := func(t reflect.Type) string {
		if t == nil {
			return "nil"
		}
		return t.Name()
	}

	for n, tc := range tt {
		t.Run(name(n), func(t *testing.T) {
			e, ok := reflector.Elem(n)
			if assert.Equal(t, tc, ok) && tc {
				assert.Equal(t, n.Elem(), e)
			}
		})
	}
}

func TestEnsurePointer(t *testing.T) {
	v := reflect.ValueOf("this is a string")
	v = reflector.EnsurePointer(v)
	strPtr := reflect.PointerTo(ltype.String)
	assert.Equal(t, strPtr, v.Type())

	v = reflector.EnsurePointer(v)
	assert.Equal(t, strPtr, v.Type())
}
