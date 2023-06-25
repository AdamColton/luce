package lusers

import (
	"sort"

	"golang.org/x/crypto/bcrypt"
)

// User stores the ID, Name and HashedPassword to validate logins. It also
// stores the Groups the user belongs to.
type User struct {
	ID             []byte `json:"-"`
	Name           string
	HashedPassword []byte
	// Groups are storted for fast searching
	Groups []string
}

// SetPassword uses bcrypt to set a HashedPassword
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = hashedPassword
	return nil
}

// CheckPassword uses bcrypt to check the password against HashedPassword.
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
}

// In checks if the user is in the group.
func (u *User) In(group string) bool {
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
