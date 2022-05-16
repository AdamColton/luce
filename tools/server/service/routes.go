package service

import (
	"fmt"
	"strings"

	"github.com/adamcolton/luce/lerr"
)

// Routes holds a collection of RouteConfigs.
type Routes []RouteConfig

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (Routes) TypeID32() uint32 {
	return 2150347030
}

// RouteConfigGen is used to create RouteConfigs from a Base path.
type RouteConfigGen struct {
	Base string
}

func NewRouteConfigGen(basePath string) *RouteConfigGen {
	return &RouteConfigGen{
		Base: basePath,
	}
}

// Get creates a RouteConfig with a Path of Base+path at
func (g *RouteConfigGen) Get(path string) *RouteConfig {
	r := &RouteConfig{
		Path:   g.Base + path,
		Method: "GET",
	}
	return r
}

// GetQuery
func (g *RouteConfigGen) GetQuery(path string) *RouteConfig {
	r := &RouteConfig{
		Path:   g.Base + path,
		Method: "GET",
		Query:  true,
	}
	return r
}

func (g *RouteConfigGen) Post(path string) *RouteConfig {
	r := &RouteConfig{
		Path:   g.Base + path,
		Method: "POST",
	}
	return r
}

func (g *RouteConfigGen) PostForm(path string) *RouteConfig {
	r := &RouteConfig{
		Path:   g.Base + path,
		Method: "POST",
		Form:   true,
	}
	return r
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
}

const ErrPathRequired = lerr.Str("RouteConfig: Path is required")

func (r *RouteConfig) WithUser() *RouteConfig {
	r.User = true
	return r
}

func (r *RouteConfig) WithPrefix() *RouteConfig {
	r.PathPrefix = true
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
		r.ID = r.String()
	}
	if strings.Contains(r.Path, "{") {
		r.PathVars = true
	}
	return nil
}

func (r *RouteConfig) String() string {
	prefxStr := ""
	if r.PathPrefix {
		prefxStr = "..."
	}
	return fmt.Sprintf("(%s) %s%s", r.Method, r.Path, prefxStr)
}

func (r *RouteConfig) Methods() []string {
	out := strings.Split(r.Method, ",")
	for i, m := range out {
		out[i] = strings.TrimSpace(m)
	}
	return out
}
