package lusers

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

func TestGroups(t *testing.T) {
	storeFac := ephemeral.Factory(bytebtree.New, 10)

	us, err := NewUserStore(storeFac)
	assert.NoError(t, err)

	expected := []string{"admin", "editor", "user"}
	for _, name := range expected {
		_, err := us.Group(name)
		assert.NoError(t, err)
	}

	got := us.Groups()
	assert.Equal(t, expected, got)

	u, err := us.Create("testUser", "password")
	assert.NoError(t, err)

	g, err := us.Group("admin")
	assert.NoError(t, err)
	g.AddUser(u)
	assert.Equal(t, "admin", u.Groups[0])
	assert.True(t, g.HasUser(u))
	g2 := us.HasGroup("admin")
	assert.Equal(t, g, g2)
	g2 = us.HasGroup("not-a-group")
	assert.Nil(t, g2)
}

func TestInGroup(t *testing.T) {
	u := &User{
		Groups: []string{"foo", "bar", "glorp"},
	}
	u.sortGroups()

	assert.True(t, u.In("foo"))
	assert.True(t, u.In("bar"))
	assert.True(t, u.In("glorp"))
	assert.False(t, u.In("baz"))
	assert.False(t, u.In("aaa"))
	assert.False(t, u.In("zzz"))
}
