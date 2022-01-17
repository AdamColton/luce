package service

import (
	"fmt"
	"strings"

	"github.com/adamcolton/luce/lerr"
)

type Require struct {
	Group string
}

func (r *Require) RequireGroup(group string) {
	if r.Group == "" {
		r.Group = group
	} else {
		r.Group += "," + group
	}
}

// Routes holds a collection of RouteConfigs.
type Routes []RouteConfig

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (Routes) TypeID32() uint32 {
	return 2150347030
}

// RouteConfigGen is used to create RouteConfigs from a Base path.
type RouteConfigGen struct {
	Base string
	Require
}

func NewRouteConfigGen(basePath string) *RouteConfigGen {
	return &RouteConfigGen{
		Base: basePath,
	}
}

func (r *RouteConfigGen) RequireGroup(group string) *RouteConfigGen {
	r.Require.RequireGroup(group)
	return r
}

func (g *RouteConfigGen) Path(path string) *RouteConfig {
	r := NewRoute(g.Base + path)
	r.Require = g.Require
	return r
}

// Get creates a RouteConfig with a Path of Base+path at
func (g *RouteConfigGen) Get(path string) *RouteConfig {
	return g.Path(path).Get()
}

// GetQuery
func (g *RouteConfigGen) GetQuery(path string) *RouteConfig {
	return g.Path(path).Get().WithQuery()
}

func (g *RouteConfigGen) Post(path string) *RouteConfig {
	return g.Path(path).Post()
}

func (g *RouteConfigGen) PostForm(path string) *RouteConfig {
	return g.Path(path).Post().WithForm()
}

type RouteConfig struct {
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

const ErrPathRequired = lerr.Str("RouteConfig: Path is required")

func NewRoute(path string) *RouteConfig {
	return &RouteConfig{
		Path: path,
	}
}

func (r *RouteConfig) AddMethod(method string) *RouteConfig {
	if r.Method == "" {
		r.Method = method
	} else {
		r.Method += "," + method
	}
	return r
}

func (r *RouteConfig) Get() *RouteConfig    { return r.AddMethod("GET") }
func (r *RouteConfig) Post() *RouteConfig   { return r.AddMethod("POST") }
func (r *RouteConfig) Delete() *RouteConfig { return r.AddMethod("DELETE") }
func (r *RouteConfig) Put() *RouteConfig    { return r.AddMethod("PUT") }

func (r *RouteConfig) WithQuery() *RouteConfig {
	r.Query = true
	return r
}

func (r *RouteConfig) WithForm() *RouteConfig {
	r.Form = true
	return r
}

func (r *RouteConfig) WithUser() *RouteConfig {
	r.User = true
	return r
}

func (r *RouteConfig) WithBody() *RouteConfig {
	r.Body = true
	return r
}

func (r *RouteConfig) WithPrefix() *RouteConfig {
	r.PathPrefix = true
	return r
}

func (r *RouteConfig) RequireGroup(group string) *RouteConfig {
	r.Require.RequireGroup(group)
	return r
}

// Validate the RouteConfig has necessary fields filled in. Unset fields will
// be set to their defaults.
func (r *RouteConfig) Validate() error {
	if r.Path == "" {
		return ErrPathRequired
	}
	if r.Method == "" {
		r.Method = "GET"
	}
	if r.ID == "" {
		r.ID = fmt.Sprintf("(%s) %s", r.Method, r.Path)
		if r.PathPrefix {
			r.ID += "..."
		}
	}
	if strings.Contains(r.Path, "{") {
		r.PathVars = true
	}
	return nil
}

func (r *RouteConfig) Methods() []string {
	out := strings.Split(r.Method, ",")
	for i, m := range out {
		out[i] = strings.TrimSpace(m)
	}
	return out
}
