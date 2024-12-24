package service

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lstr"
)

// Require breaks out the possible route requirements so the logic can be shared
// between RouteConfigGen and RouteConfig.
type Require struct {
	Group string
}

// RequireGroup sets or appends the group to the Group field.
func (r *Require) RequireGroup(group string) {
	if r.Group == "" {
		r.Group = group
	} else {
		r.Group += "," + group
	}
}

var slash = lstr.Seperator("/")

const ErrPathRequired = lerr.Str("RouteConfig: Path is required")
