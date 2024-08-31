package entity_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/lerr"
	lgob "github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/entity"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

type Sub struct {
	id   byte
	Name string
}

func (s *Sub) EntKey() []byte {
	return []byte{s.id}
}

type Core struct {
	id        []byte
	Primary   entity.Reference[*Sub]
	Secondary entity.Reference[*Sub]
}

func (c *Core) EntKey() []byte {
	return c.id
}

func TestReferenceAndBuilder(t *testing.T) {
	builders := map[string]func(store.Factory) entity.Builder{
		"json": entity.NewJsonBuilder,
		"gob:": entity.NewGobBuilder,
	}

	for n, builderFn := range builders {
		t.Run(n, func(t *testing.T) {
			s, err := ephemeral.Factory(bytebtree.New, 1).Store([]byte("test"))
			assert.NoError(t, err)
			b := builderFn(s)
			coreStore := lerr.Must(entity.NewStore[*Core](b, "core", nil))
			subStore := lerr.Must(entity.NewStore[*Sub](b, "sub", nil))

			entity.AddGetter(subStore)

			c := &Core{
				id: []byte{1, 2},
				Primary: entity.NewRef(&Sub{
					id:   4,
					Name: "Adam",
				}),
				Secondary: entity.NewRef(&Sub{
					id:   55,
					Name: "Lauren",
				}),
			}

			assert.Equal(t, "Adam", c.Primary.Ent.Name)

			_, err = coreStore.Put(c, nil)
			assert.NoError(t, err)
			_, err = subStore.Put(c.Primary.Ent, nil)
			assert.NoError(t, err)
			_, err = subStore.Put(c.Secondary.Ent, nil)
			assert.NoError(t, err)

			found, c2, err := coreStore.Get(c.id)
			assert.True(t, found)
			assert.NoError(t, err)
			assert.Nil(t, c2.Primary.Ent)

			p, ok := c2.Primary.Get()
			assert.True(t, ok)
			assert.Equal(t, "Adam", p.Name)
			assert.Equal(t, "Adam", c.Primary.Ent.Name)
		})
	}
}

func TestRefGob(t *testing.T) {
	s := &Sub{
		id:   123,
		Name: "gob-test",
	}
	r := entity.NewRef(s)

	buf := lgob.Enc(&r)

	r2 := entity.Reference[*Sub]{}
	lgob.Dec(buf, &r2)

	assert.Equal(t, r.ID, r2.ID)
}

func TestRefPtrGob(t *testing.T) {
	s := &Sub{
		id:   123,
		Name: "gob-test",
	}
	r := entity.NewRef(s)
	var ptr graph.Ptr[*Sub] = &r

	buf := lgob.Enc(ptr)
	assert.NotNil(t, buf)

	ptr2 := entity.Reference[*Sub]{}
	lgob.Dec(buf, &ptr2)
	assert.NotNil(t, ptr2.ID)
}
