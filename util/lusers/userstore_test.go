package lusers

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

func TestUserStore(t *testing.T) {
	storeFac := ephemeral.Factory(bytebtree.New, 10)

	us := MustUserStore(storeFac)

	names := []string{"user1", "user2", "user3", "user4"}
	for _, n := range names {
		pwd := n + "-password"
		uc, err := us.Create(n, pwd)
		assert.NoError(t, err)
		ul, err := us.Login(n, pwd)
		assert.NoError(t, err)
		assert.Equal(t, uc, ul)
		ul, err = us.Login(n, "bad-password")
		assert.Error(t, err)
		assert.Nil(t, ul)
	}
	ul, err := us.Login("bad-user", "bad-password")
	assert.Error(t, err)
	assert.Nil(t, ul)

	found := make(map[string]bool)
	for _, name := range us.List() {
		found[name] = true
	}

	assert.Len(t, found, len(names))
	for _, n := range names {
		assert.True(t, found[n])
	}

	_, err = us.Create("user1", "user1-password")
	assert.Equal(t, ErrUserAlreadyExists("user1"), err)
	assert.Equal(t, "User user1 already exists", err.Error())

}
