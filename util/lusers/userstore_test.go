package lusers

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store/ephemeral/quicknested"
	"github.com/stretchr/testify/assert"
)

func TestUserStore(t *testing.T) {
	us := lerr.Must(NewUserStore(quicknested.New(10)))

	names := []string{"user1", "user2", "user3", "user4"}
	for _, n := range names {
		_, err := us.Create(n, n+"-password")
		assert.NoError(t, err)
	}

	found := make(map[string]bool)
	for _, name := range us.List() {
		found[name] = true
	}

	assert.Len(t, found, len(names))
	for _, n := range names {
		assert.True(t, found[n])
	}

	_, err := us.Create("user1", "user1-password")
	assert.Equal(t, ErrUserAlreadyExists("user1"), err)
	assert.Equal(t, "User user1 already exists", err.Error())
}
