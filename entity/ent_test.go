package entity_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/entity/enttest"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	enttest.Setup()
	entity.ClearCache()

	id := rye.Serialize.Any(31415, nil)
	foo := &enttest.Foo{
		ID:   id,
		Name: "Adam Colton",
	}

	ref := entity.Put(foo)
	f2, ok := ref.Get()
	assert.True(t, ok)
	assert.Equal(t, "Adam Colton", f2.Name)

	// Check that both are pointing to the same instance
	assert.Equal(t, foo, f2)

	err := enttest.SaveAndWait(ref, false)
	assert.NoError(t, err)

	// calling clear now does clear the cache
	// ...and I think that's right
	// ref.Clear(false)
	// f2, ok = ref.Get()
	// foo.Name = "UPDATE VIA CACHE"
	// assert.True(t, ok)
	// assert.Equal(t, "UPDATE VIA CACHE", f2.Name)

	ref.Clear(true)
	f2, ok = ref.WeakGet()
	assert.False(t, ok)
	assert.Nil(t, f2)

	f2, ok = ref.Get()
	foo.Name = "NOT REFERENCED"
	assert.True(t, ok)
	assert.Equal(t, "Adam Colton", f2.Name)

	ref2, ok := entity.Get[enttest.Foo](id)
	assert.True(t, ok)
	f3, ok := ref2.Get()
	assert.True(t, ok)
	assert.Equal(t, "Adam Colton", f3.Name)
}

func TestGarbage(t *testing.T) {
	entity.ClearCache()
	m32 := type32.NewTypeMap()
	entity.Setup{
		Store:        lerr.Must(ephemeral.Factory(bytebtree.New, 1).FlatStore([]byte("testing"))),
		Typer:        m32,
		Serializer:   gob.Serializer{},
		Deserializer: gob.Deserializer{},
	}.Init()
	err := serial.RegisterPtr[enttest.Foo](m32)
	assert.NoError(t, err)

	root := &enttest.Foo{
		ID:   entity.Rand(),
		Name: "the root",
	}

	hasRef := &enttest.Foo{
		ID:   entity.Rand(),
		Name: "has ref",
	}
	root.Refs = append(root.Refs, hasRef.ID)

	missingRoot := &enttest.Foo{
		ID:   entity.Rand(),
		Name: "the garbage",
	}

	entity.Put(root).Save(root)
	entity.Put(hasRef).Save(hasRef)
	entity.Put(missingRoot).Save(missingRoot)
	entity.AddGCRoots(root.EntKey())

	g := entity.Garbage()
	assert.NoError(t, err)
	// garbage always returns nothing on first pass
	assert.Len(t, g, 0)

	g = entity.Garbage()
	assert.Equal(t, g[0], missingRoot.EntKey())
}

func TestInterface(t *testing.T) {
	enttest.Setup()
	entity.ClearCache()

	s := &enttest.String{
		Key32: entity.Rand32(),
		Str:   "this is a test",
	}
	sr := entity.Put(s)
	err := enttest.SaveAndWait(sr, false)
	assert.NoError(t, err)

	i := &enttest.Int{
		Key32: entity.Rand32(),
		I:     31415,
	}
	ir := entity.Put(i)
	err = enttest.SaveAndWait(ir, false)
	assert.NoError(t, err)

	e, err := entity.Load(ir.EntKey())
	assert.NoError(t, err)
	assert.Equal(t, i, e)

	e, err = entity.Load(sr.EntKey())
	assert.NoError(t, err)
	assert.Equal(t, s, e)
}
