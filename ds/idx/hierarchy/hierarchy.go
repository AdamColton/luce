package hierarchy

import (
	"github.com/adamcolton/luce/ds/bimap"
	"github.com/adamcolton/luce/ds/lset"
	"golang.org/x/exp/constraints"
)

type Hierarchy[ID constraints.Integer, Name comparable] struct {
	*bimap.Bimap[Key[ID, Name], ID]
	Children map[ID]*lset.Set[Name]
	MaxID    ID
}

func New[ID constraints.Integer, Name comparable](size int) *Hierarchy[ID, Name] {
	return &Hierarchy[ID, Name]{
		Bimap: bimap.New[Key[ID, Name], ID](size),
		Children: map[ID]*lset.Set[Name]{
			0: lset.New[Name](),
		},
		MaxID: 1, // reserve 0 for root
	}
}

type Key[ID, Name comparable] struct {
	ID   ID
	Name Name
}

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
