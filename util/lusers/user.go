package lusers

import (
	"sort"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/filter"
	"golang.org/x/crypto/bcrypt"
)

// User stores the ID, Name and HashedPassword to validate logins. It also
// stores the Groups the user belongs to.
type User struct {
	ID             []byte `json:"-"`
	Name           string
	HashedPassword []byte
	// Groups are storted for fast searching
	Groups slice.Slice[string]
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
	if u == nil {
		return false
	}
	idx := u.gidx(group)
	return u.Groups.IdxCheck(idx) && u.Groups[idx] == group
}

func (u *User) gidx(group string) int {
	return u.Groups.Search(filter.GTE(group))
}

func (u *User) sortGroups() {
	sort.Slice(u.Groups, func(i, j int) bool {
		return u.Groups[i] < u.Groups[j]
	})
}

// OneRequired requires that a user be in one of the listed groups. If no groups
// are listed, the result is true.
func (u *User) OneRequired(groups []string) bool {
	if len(groups) == 0 {
		return true
	}
	if u == nil {
		return false
	}

	var group string
	fn := func(g string) bool {
		return g >= group
	}

	for _, group = range groups {
		idx := u.Groups.Search(fn)
		if u.Groups.IdxCheck(idx) && u.Groups[idx] == group {
			return true
		}
	}
	return false
}
