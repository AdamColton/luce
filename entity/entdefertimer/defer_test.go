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

func TestEntity(t *testing.T) {
	entity.ClearCache()
	entity.DeferStrategy = entdefertimer.New(time.Millisecond*10, time.Millisecond)
	enttest.Setup()

	id := rye.Serialize.Any(31415, nil)
	foo := &enttest.Foo{
		ID:   id,
		Name: "init",
	}

	ref := entity.Put(foo)
	ref.Save(foo)
	ref.Save(foo) // should cause a reset
	foo.Name = "after save invoked, before triggered"
	err := enttest.Wait(ref.EntKey())
	assert.NoError(t, err)

	ref.Clear(true)
	f2, ok := ref.WeakGet()
	assert.False(t, ok)
	assert.Nil(t, f2)
	f2, ok = ref.Get()
	assert.True(t, ok)
	assert.Equal(t, "after save invoked, before triggered", f2.Name)
}
