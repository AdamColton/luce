package service

import (
	"fmt"
	"path"
	"strings"

	"github.com/adamcolton/luce/lerr"
)

type Link struct {
	Name       string
	Host, Path string
}

func (l Link) Get(host, port string) string {
	if l.Host == "" {
		return l.Path
	}
	return fmt.Sprintf("https://%s.%s%s%s", l.Host, host, port, l.Path)
}

type Service struct {
	Name   string
	Host   string
	Base   string
	Routes []ServiceRoute
	Links  []Link
}

func (*Service) TypeID32() uint32 {
	return 2516527266
}

func (s *Service) Validate() error {
	return lerr.NewSliceErrs(len(s.Routes), -1, func(i int) error {
		r := &(s.Routes[i])
		return r.Validate()
	})
}

func (s *Service) AddLink(name, host string, pth ...string) {
	wBase := make([]string, len(pth)+1)
	wBase[0] = s.Base
	copy(wBase[1:], pth)
	s.Links = append(s.Links, Link{
		Name: name,
		Host: host,
		Path: path.Join(wBase...),
	})
}

type ServiceRoute struct {
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

func NewServiceRoute(path string) *ServiceRoute {
	if !strings.HasPrefix(path, string(slash)) {
		path = string(slash) + path
	}
	return &ServiceRoute{
		Path: path,
	}
}

// AddMethod is a chainable helper. It sets the Method field. If the Method
// field is not empty, it appends the input with comma seperation.
func (r *ServiceRoute) AddMethod(method string) *ServiceRoute {
	if r.Method == "" {
		r.Method = method
	} else {
		r.Method += "," + method
	}
	return r
}

// Get is a chainable helper. It adds Get to the Method field.
func (r *ServiceRoute) Get() *ServiceRoute { return r.AddMethod("GET") }

// Post is a chainable helper. It adds Post to the Method field.
func (r *ServiceRoute) Post() *ServiceRoute { return r.AddMethod("POST") }

// Delete is a chainable helper. It adds Delete to the Method field.
func (r *ServiceRoute) Delete() *ServiceRoute { return r.AddMethod("DELETE") }

// Put is a chainable helper. It adds Put to the Method field.
func (r *ServiceRoute) Put() *ServiceRoute { return r.AddMethod("PUT") }

// WithQuery is a chainable helper. It sets the Query field to true.
func (r *ServiceRoute) WithQuery() *ServiceRoute {
	r.Query = true
	return r
}

// WithQuery is a chainable helper. It sets the Form field to true.
func (r *ServiceRoute) WithForm() *ServiceRoute {
	r.Form = true
	return r
}

// WithUser is a chainable helper. It sets the User field to true.
func (r *ServiceRoute) WithUser() *ServiceRoute {
	r.User = true
	return r
}

// WithBody is a chainable helper. It sets the Body field to true.
func (r *ServiceRoute) WithBody() *ServiceRoute {
	r.Body = true
	return r
}

// WithPrefix is a chainable helper. It sets the PathPrefix field to true.
func (r *ServiceRoute) WithPrefix() *ServiceRoute {
	r.PathPrefix = true
	return r
}

// RequireGroup is a chainable helper. It calls RequireGroup on the embedded
// Require.
func (r *ServiceRoute) RequireGroup(group string) *ServiceRoute {
	r.Require.RequireGroup(group)
	return r
}

// Validate the ServiceRoute has necessary fields filled in. Unset fields will
// be set to their defaults.
func (r *ServiceRoute) Validate() error {
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
func (r *ServiceRoute) String() string {
	prefxStr := ""
	if r.PathPrefix {
		prefxStr = "..."
	}
	return fmt.Sprintf("(%s) %s%s", r.Method, r.Path, prefxStr)
}

func (r *ServiceRoute) Methods() []string {
	out := strings.Split(r.Method, ",")
	for i, m := range out {
		out[i] = strings.TrimSpace(m)
	}
	return out
}

func (r *ServiceRoute) Handle(c *Client, h RequestResponder) {
	c.AddServiceRoute(h, r)
}
