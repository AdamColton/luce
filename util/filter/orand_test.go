package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestOrAnd(t *testing.T) {
	taskTags := map[string]bool{"work": true, "urgent": true}
	hasTag := filter.New(func(tag string) bool { return taskTags[tag] })

	// (work AND urgent) OR home
	assert.True(t, hasTag.OrAnd(filter.OrAnd[string]{{"work", "urgent"}, {"home"}}))
	// (home AND errands) OR work: the second set matches
	assert.True(t, hasTag.OrAnd(filter.OrAnd[string]{{"home", "errands"}, {"work"}}))
	// (work AND home) OR errands
	assert.False(t, hasTag.OrAnd(filter.OrAnd[string]{{"work", "home"}, {"errands"}}))

	assert.True(t, hasTag.OrAnd(nil), "no sets, nothing to satisfy")
}
