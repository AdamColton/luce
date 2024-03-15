package cmprtest

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/util/reflector"
)

type tHelper interface {
	Helper()
}

type TestingT interface {
	Errorf(format string, args ...interface{})
}

const (
	// Small is the value that will be passed into AssertEqualizer
	Small cmpr.Tolerance = 1e-10
)

// Equal calls AssertEqual with the default value of Small. If there is an error
// it is passed into t.Error. The return bool will be true if the values were
// equal.
func Equal(t TestingT, expected, actual interface{}, msg ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return EqualInDelta(t, expected, actual, Small, msg...)
}

// EqualInDelta calls AssertEqual. If there is an error it is passed into
// t.Error. The return bool will be true if the values were equal.
func EqualInDelta(t TestingT, expected, actual interface{}, delta cmpr.Tolerance, msg ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	err := AssertEqual(expected, actual, delta)
	if err == nil {
		return true
	}
	if len(msg) > 0 {
		t.Errorf("%s: %s", err.Error(), Message(msg...))
	} else {
		t.Errorf("%s", err.Error())
	}
	return false
}

var equalType = reflector.Type[cmpr.AssertEqualizer]()

// AssertEqual can compare anything that implements geomtest.AssertEqualizer.
// There is also logic to handle comparing float64 values Any two slices whose
// elements can be compared with Equal can be compared. The provided delta value
// will be passed to anything that implements AssertEqualizer. If the equality
// check fails, an error is returned.
func AssertEqual(expected, actual interface{}, delta cmpr.Tolerance) error {
	ev := reflect.ValueOf(expected)

	if eq, ok := expected.(cmpr.AssertEqualizer); ok {
		return eq.AssertEqual(actual, delta)
	} else if ev.Kind() == reflect.Slice {
		av := reflect.ValueOf(actual)
		if av.Kind() != reflect.Slice {
			return lerr.TypeMismatch(expected, actual)
		}
		return lerr.NewSliceErrs(ev.Len(), av.Len(), func(i int) error {
			return AssertEqual(ev.Index(i).Interface(), av.Index(i).Interface(), delta)
		})
	} else if ef, ok := f64(expected); ok {
		af, ok := f64(actual)
		if !ok {
			return lerr.NewTypeMismatch(expected, actual)
		}
		return lerr.NewNotEqual(delta.Equal(ef, af), ef, af)
	}

	format := "unsupported_type: %s"
	t := ev.Type()
	if t.Kind() != reflect.Ptr {
		if p := reflect.PtrTo(t); p.Implements(equalType) {
			format = fmt.Sprintf("%s (%s fulfills AssertEqualizer)", format, p.String())
		}
	}

	return fmt.Errorf(format, t.String())
}

func f64(i any) (float64, bool) {
	switch i := i.(type) {
	case float64:
		return i, true
	case float32:
		return float64(i), true
	case int:
		return float64(i), true
	case int8:
		return float64(i), true
	case int16:
		return float64(i), true
	case int32:
		return float64(i), true
	case int64:
		return float64(i), true
	case uint:
		return float64(i), true
	case uint8:
		return float64(i), true
	case uint16:
		return float64(i), true
	case uint32:
		return float64(i), true
	case uint64:
		return float64(i), true
	}
	return 0, false
}

// Message takes in args to form a message. If there are more than 1 arg and the
// first is a string, it will use that as a format string.
func Message(msg ...interface{}) string {
	// TODO: Move this to lstr
	ln := len(msg)
	if ln == 0 {
		return ""
	}
	if ln == 1 {
		if s, ok := msg[0].(string); ok {
			return s
		}
		return fmt.Sprint(msg[0])
	}
	if f, ok := msg[0].(string); ok {
		return fmt.Sprintf(f, msg[1:]...)
	}
	return fmt.Sprint(msg...)
}
