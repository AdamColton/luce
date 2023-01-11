package list_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

type mockList struct{}

func (m mockList) AtIdx(idx int) int { return idx }
func (m mockList) Len() int          { return 10 }
func (m mockList) String() string    { return "0...9" }

func TestWrap(t *testing.T) {
	var m mockList
	w := list.Wrap[int](m)
	assert.Equal(t, m, w.List)

	w = list.Wrap[int](w)
	assert.Equal(t, m, w.List)
	_, shouldBeFalse := w.List.(list.Wrapper[int])
	assert.False(t, shouldBeFalse)
}

func TestUpgrade(t *testing.T) {
	var m mockList
	w := list.Wrap[int](m)
	var s fmt.Stringer
	assert.True(t, upgrade.Upgrade(w, &s))
	assert.Equal(t, "0...9", s.String())
}
