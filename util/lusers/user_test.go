package lusers

import (
	"testing"

	"github.com/adamcolton/luce/store/ephemeral/quicknested"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	us, err := NewUserStore(quicknested.New(10))
	assert.NoError(t, err)

	n, p := "testuser", "password"
	u, err := us.Create(n, p)
	assert.NoError(t, err)

	u2, err := us.GetByID(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u, u2)
	assert.Equal(t, u.HashedPassword, u2.HashedPassword)
	assert.NoError(t, u2.CheckPassword(p))
	assert.Equal(t, n, u2.Name)
}
