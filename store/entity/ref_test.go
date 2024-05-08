package entity_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/lerr"
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

func TestReference(t *testing.T) {
	s, err := ephemeral.Factory(bytebtree.New, 1).Store([]byte("test"))
	assert.NoError(t, err)
	b := entity.NewJsonBuilder(s)
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
}
