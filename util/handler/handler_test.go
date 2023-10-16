package handler_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tt := []struct {
		fn     any
		hasRet bool
		errIdx int
		err    error
		reflect.Type
		arg, out any
	}{
		{
			fn:     func() {},
			hasRet: false,
			errIdx: -1,
		}, {
			fn:     func(string) {},
			hasRet: false,
			errIdx: -1,
			Type:   ltype.String,
			arg:    "test",
		}, {
			fn:     func() string { return "test" },
			hasRet: true,
			errIdx: -1,
			out:    "test",
		}, {
			fn:     func(string) string { return "test" },
			hasRet: true,
			errIdx: -1,
			Type:   ltype.String,
			arg:    "test",
			out:    "test",
		}, {
			fn:     func() error { return nil },
			hasRet: false,
			errIdx: 0,
		}, {
			fn:     func(string) error { return nil },
			hasRet: false,
			errIdx: 0,
			Type:   ltype.String,
			arg:    "test",
		}, {
			fn:     func() (string, error) { return "test", nil },
			hasRet: true,
			errIdx: 1,
			out:    "test",
		}, {
			fn:     func(string) (string, error) { return "test", nil },
			hasRet: true,
			errIdx: 1,
			Type:   ltype.String,
			arg:    "test",
			out:    "test",
		}, {
			fn:     func() error { return lerr.Str("test error") },
			errIdx: 0,
			err:    lerr.Str("test error"),
		}, {
			fn:  "test",
			err: lerr.Str("expected func, got: string"),
		}, {
			fn:  func() (string, string) { return "test", "test" },
			err: lerr.Str("expected func(T?) (U?, error?) where ? indicates optional, got: func() (string, string)"),
		}, {
			fn:  func(string, string) {},
			err: lerr.Str("expected func(T?) (U?, error?) where ? indicates optional, got: func(string, string)"),
		}, {
			fn:  func() (error, string) { return nil, "test" },
			err: lerr.Str("expected func(T?) (U?, error?) where ? indicates optional, got: func() (error, string)"),
		},
	}

	for _, tc := range tt {
		n := strings.ReplaceAll(reflect.TypeOf(tc.fn).String(), " ", "_")
		t.Run(n, func(t *testing.T) {
			h, err := handler.New(tc.fn)
			if h != nil {
				assert.Equal(t, tc.hasRet, h.HasRet())
				assert.Equal(t, tc.errIdx, h.ErrIdx())
				assert.Equal(t, tc.Type, h.Type())
				out, err := h.Handle(tc.arg)
				assert.Equal(t, tc.out, out)
				assert.Equal(t, tc.err, err)
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
