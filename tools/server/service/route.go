package service

import (
	"fmt"
	"strings"

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

type Route struct {
	ID   string
	Path string
	// Method is comma delimited methods that are accepted
	Method     string
	PathPrefix bool
	PathVars   bool
	Form       bool
	Body       bool
	User       bool
	Query      bool
	Require
}

func NewRoute(path string) *Route {
	if !strings.HasPrefix(path, string(slash)) {
		path = string(slash) + path
	}
	return &Route{
		Path: path,
	}
}

// AddMethod is a chainable helper. It sets the Method field. If the Method
// field is not empty, it appends the input with comma seperation.
func (r *Route) AddMethod(method string) *Route {
	if r.Method == "" {
		r.Method = method
	} else {
		r.Method += "," + method
	}
	return r
}

// Get is a chainable helper. It adds Get to the Method field.
func (r *Route) Get() *Route { return r.AddMethod("GET") }

// Post is a chainable helper. It adds Post to the Method field.
func (r *Route) Post() *Route { return r.AddMethod("POST") }

// Delete is a chainable helper. It adds Delete to the Method field.
func (r *Route) Delete() *Route { return r.AddMethod("DELETE") }

// Put is a chainable helper. It adds Put to the Method field.
func (r *Route) Put() *Route { return r.AddMethod("PUT") }

// WithQuery is a chainable helper. It sets the Query field to true.
func (r *Route) WithQuery() *Route {
	r.Query = true
	return r
}

// WithQuery is a chainable helper. It sets the Form field to true.
func (r *Route) WithForm() *Route {
	r.Form = true
	return r
}

// WithUser is a chainable helper. It sets the User field to true.
func (r *Route) WithUser() *Route {
	r.User = true
	return r
}

// WithBody is a chainable helper. It sets the Body field to true.
func (r *Route) WithBody() *Route {
	r.Body = true
	return r
}

// WithPrefix is a chainable helper. It sets the PathPrefix field to true.
func (r *Route) WithPrefix() *Route {
	r.PathPrefix = true
	return r
}

// RequireGroup is a chainable helper. It calls RequireGroup on the embedded
// Require.
func (r *Route) RequireGroup(group string) *Route {
	r.Require.RequireGroup(group)
	return r
}

// Validate the ServiceRoute has necessary fields filled in. Unset fields will
// be set to their defaults.
func (r *Route) Validate() error {
	if r.Path == "" {
		return ErrPathRequired
	}
	if r.Method == "" {
		r.Method = "GET"
	}
	if r.ID == "" {
		r.ID = r.String()
	}
	if strings.Contains(r.Path, "{") {
		r.PathVars = true
	}
	return nil
}

// String fulfills Stringer. Returns the ServiceRoute as "(Method)path...", where
// the ellipse will only be present if the path is a prefix.
func (r *Route) String() string {
	prefxStr := ""
	if r.PathPrefix {
		prefxStr = "..."
	}
	return fmt.Sprintf("(%s) %s%s", r.Method, r.Path, prefxStr)
}

func (r *Route) Methods() []string {
	out := strings.Split(r.Method, ",")
	for i, m := range out {
		out[i] = strings.TrimSpace(m)
	}
	return out
}
