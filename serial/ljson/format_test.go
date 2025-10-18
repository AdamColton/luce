package ljson_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestFormatWriteTo(t *testing.T) {
	type Person struct {
		Name    string
		Tags    []string
		Address struct{ City string }
	}
	p := Person{Name: "Adam", Tags: []string{"a", "b"}}
	p.Address.City = "Boston"

	ctx := ljson.NewMarshalContext("")
	ctx.Sort = true
	wn := lerr.Must(ljson.Marshal(p, ctx))

	buf := bytes.NewBuffer(nil)
	_, err := wn.FormatWriteTo(buf, "\n", "  ", false)
	assert.NoError(t, err)
	assert.Equal(t, "{\n  \"Address\":{\n    \"City\":\"Boston\"\n  },\n  \"Name\":\"Adam\",\n  \"Tags\":[\"a\",\"b\"]\n}", buf.String())

	buf.Reset()
	_, err = wn.FormatWriteTo(buf, "", "", false)
	assert.NoError(t, err)
	assert.Equal(t, wn.String(), buf.String(), "without a prefix and indent it is compact")
}
