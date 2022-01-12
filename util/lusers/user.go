package lusers

import (
	"sort"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             []byte `json:"-"`
	Name           string
	HashedPassword []byte
	Groups         []string
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = hashedPassword
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
}

func (u *User) In(group string) bool {
	if u == nil {
		return false
	}
	idx := sort.Search(len(u.Groups), func(i int) bool {
		return u.Groups[i] >= group
	})
	if idx < 0 || idx >= len(u.Groups) {
		return false
	}
	return u.Groups[idx] == group
}

func (u *User) sortGroups() {
	sort.Slice(u.Groups, func(i, j int) bool {
		return u.Groups[i] < u.Groups[j]
	})
}
