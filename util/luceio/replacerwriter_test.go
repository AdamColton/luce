package luceio_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestReplacerWriter(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	r := luceio.NewReplacer(buf, "foo", "bar", "3", "5")
	r.WriteString("test 3 test foo 3foo3")
	assert.Equal(t, "test 5 test bar 5bar5", buf.String())
}
