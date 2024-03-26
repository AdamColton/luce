// Package hierarchy maps paths of names to integer IDs, with each name
// registered under its parent.
package hierarchy

import (
	"github.com/adamcolton/luce/ds/bimap"
	"github.com/adamcolton/luce/ds/lset"
	"golang.org/x/exp/constraints"
)

// Hierarchy maps a path of names, like a path in a file system, to an ID. Each
// name is registered under its parent, so the same name in two places has two
// IDs. The root has the ID 0. It is not safe for concurrent use.
type Hierarchy[ID constraints.Integer, Name comparable] struct {
	// Bimap maps the Key of each node to its ID. Adding to it directly leaves the
	// Hierarchy inconsistent, use Get and Key.
	*bimap.Bimap[Key[ID, Name], ID]
	// Children holds the names registered under each ID.
	Children map[ID]*lset.Set[Name]
	// MaxID is the next ID to be assigned. IDs are not reused, and nothing checks
	// that MaxID fits in ID.
	MaxID ID
}

// New creates a Hierarchy with room for size names. It only has the root.
func New[ID constraints.Integer, Name comparable](size int) *Hierarchy[ID, Name] {
	return &Hierarchy[ID, Name]{
		Bimap: bimap.New[Key[ID, Name], ID](size),
		Children: map[ID]*lset.Set[Name]{
			0: lset.New[Name](),
		},
		MaxID: 1, // reserve 0 for root
	}
}

// Key identifies a name by the ID of its parent.
type Key[ID, Name comparable] struct {
	ID   ID
	Name Name
}

// Get returns the ID of the last name in the path, following it from the root.
// An empty path is the root, 0. If create is true, names that are not there are
// added. The bool is true if the last name was already there, so it is false
// when it was created. If create is false and a name is missing, it returns 0
// and false.
func (h *Hierarchy[ID, Name]) Get(path []Name, create bool) (id ID, found bool) {
	if len(path) == 0 {
		return 0, true
	}
	for _, name := range path {
		id, found = h.Key(id, name, create)
		if !found && !create {
			return
		}
	}
	return
}

// Key returns the ID of the name under the parent id, and true if it was
// already there. If it is not there and create is true, it is added under the
// next ID and that ID is returned with false; if create is false it returns 0
// and false. It panics if create is true and id is not the ID of a name in the
// Hierarchy.
func (h *Hierarchy[ID, Name]) Key(id ID, name Name, create bool) (ID, bool) {
	k := Key[ID, Name]{
		ID:   id,
		Name: name,
	}
	out, found := h.A(k)
	if !found && create {
		out = h.MaxID
		h.MaxID++
		h.Bimap.Add(k, out)
		h.Children[out] = lset.New[Name]()
		h.Children[k.ID].Add(k.Name)
	}
	return out, found
}
