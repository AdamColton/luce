package entdefertimer_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/entity/entdefertimer"
	"github.com/adamcolton/luce/entity/enttest"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestTOQ(t *testing.T) {
	entity.ClearCache()
	saveDelay := time.Millisecond * 10
	clearDelay := saveDelay * 10
	toq := entdefertimer.NewToq(clearDelay, saveDelay, 100)
	entity.DeferStrategy = toq
	enttest.Setup()

	id := rye.Serialize.Any(31415, nil)
	foo := &enttest.Foo{
		ID:   id,
		Name: "init",
	}

	ref := entity.Put(foo)
	ref.Save(foo)
	ref.Save(foo) // should cause a reset
	_, ok := ref.WeakGet()
	assert.True(t, ok)

	time.Sleep(saveDelay * 2)
	assert.True(t, toq.Done())

	time.Sleep(clearDelay * 2)
	_, ok = ref.WeakGet()
	assert.False(t, ok)
}
