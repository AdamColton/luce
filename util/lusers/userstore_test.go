package lusers

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
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

func TestNewUserStoreErrors(t *testing.T) {
	_, err := NewUserStore(failFactory{NestedFactory: quicknested.New(10), failOn: byID})
	assert.Error(t, err)

	_, err = NewUserStore(failFactory{NestedFactory: quicknested.New(10), failOn: byName})
	assert.Error(t, err)

	_, err = NewUserStore(failFactory{NestedFactory: quicknested.New(10), failOn: groups})
	assert.Error(t, err)
}

func TestCreateErrors(t *testing.T) {
	us := lerr.Must(NewUserStore(quicknested.New(10)))

	// GetByName succeeds (the name is taken) but the stored record is
	// corrupt, so GetByID's json.Unmarshal fails with something other
	// than ErrUserNotFound.
	assert.NoError(t, us.byName.Put([]byte("broken"), []byte("bad-id")))
	assert.NoError(t, us.byID.Put([]byte("bad-id"), []byte("not json")))
	_, err := us.Create("broken", "password")
	assert.Error(t, err)
	assert.NotEqual(t, ErrUserNotFound, err)

	// SetPassword fails: bcrypt rejects passwords over 72 bytes.
	longPassword := string(make([]byte, 73))
	_, err = us.Create("user1", longPassword)
	assert.Error(t, err)

	// Update (us.byID.Put) fails.
	us2 := lerr.Must(NewUserStore(failFactory{
		NestedFactory: quicknested.New(10),
		wrapOn:        byID,
		wrap: func(s store.NestedStore) store.NestedStore {
			return errPutStore{s}
		},
	}))
	_, err = us2.Create("user1", "password")
	assert.Error(t, err)

	// us.byName.Put fails.
	us3 := lerr.Must(NewUserStore(failFactory{
		NestedFactory: quicknested.New(10),
		wrapOn:        byName,
		wrap: func(s store.NestedStore) store.NestedStore {
			return errPutStore{s}
		},
	}))
	_, err = us3.Create("user1", "password")
	assert.Error(t, err)
}

func TestGroupError(t *testing.T) {
	us := lerr.Must(NewUserStore(failFactory{
		NestedFactory: quicknested.New(10),
		wrapOn:        groups,
		wrap: func(s store.NestedStore) store.NestedStore {
			return errNestStore{s}
		},
	}))
	_, err := us.Group("admin")
	assert.Error(t, err)
}
