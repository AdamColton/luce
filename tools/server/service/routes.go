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

// Routes holds a collection of RouteConfigs.
type Routes []RouteConfig

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (Routes) TypeID32() uint32 {
	return 2150347030
}

// RouteConfigGen is used to create RouteConfigs from a Base path.
type RouteConfigGen struct {
	Base string
	Host string
	User bool
	Require
}

func NewRouteConfigGen(basePath string) *RouteConfigGen {
	return &RouteConfigGen{
		Base: basePath,
	}
}

var slash = lstr.Seperator("/")

// RequireGroup is a chainable helper. It calls RequireGroup on the embedded
// Require.
func (r *RouteConfigGen) RequireGroup(group string) *RouteConfigGen {
	r.Require.RequireGroup(group)
	return r
}

// SetHost is a chainable helper. It sets the Host field.
func (r *RouteConfigGen) SetHost(host string) *RouteConfigGen {
	r.Host = host
	return r
}

func (g *RouteConfigGen) WithUser() *RouteConfigGen {
	g.User = true
	return g
}

// Path creates a RouteConfig appending path to the Base and copying the Require
// object.
func (g *RouteConfigGen) Path(path string) *RouteConfig {
	r := NewRoute(slash.Join(g.Base, path))
	r.Require = g.Require
	r.Host = g.Host
	r.User = g.User
	return r
}

// Get creates a RouteConfig with a Path of Base+path at
func (g *RouteConfigGen) Get(path string) *RouteConfig {
	return g.Path(path).Get()
}

// GetQuery creates a clone of the RouteConfigGen appending path to the Base and
// setting the method to Get.
func (g *RouteConfigGen) GetQuery(path string) *RouteConfig {
	return g.Path(path).Get().WithQuery()
}

// GetQuery creates a clone of the RouteConfigGen appending path to the Base and
// setting the method to Post.
func (g *RouteConfigGen) Post(path string) *RouteConfig {
	return g.Path(path).Post()
}

// GetQuery creates a clone of the RouteConfigGen appending path to the Base and
// setting the method to Post and enabling form data to be included in the
// requets.
func (g *RouteConfigGen) PostForm(path string) *RouteConfig {
	return g.Path(path).Post().WithForm()
}

type RouteConfig struct {
	ID   string
	Path string
	Host string
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

// NewRoute defined by path (relative to the base path).
func NewRoute(path string) *RouteConfig {
	return &RouteConfig{
		Path: path,
	}
}

// SetHost is a chainable helper. It sets the Host field.
func (r *RouteConfig) SetHost(host string) *RouteConfig {
	r.Host = host
	return r
}

// AddMethod is a chainable helper. It sets the Method field. If the Method
// field is not empty, it appends the input with comma seperation.
func (r *RouteConfig) AddMethod(method string) *RouteConfig {
	if r.Method == "" {
		r.Method = method
	} else {
		r.Method += "," + method
	}
	return r
}

// Get is a chainable helper. It adds Get to the Method field.
func (r *RouteConfig) Get() *RouteConfig { return r.AddMethod("GET") }

// Post is a chainable helper. It adds Post to the Method field.
func (r *RouteConfig) Post() *RouteConfig { return r.AddMethod("POST") }

// Delete is a chainable helper. It adds Delete to the Method field.
func (r *RouteConfig) Delete() *RouteConfig { return r.AddMethod("DELETE") }

// Put is a chainable helper. It adds Put to the Method field.
func (r *RouteConfig) Put() *RouteConfig { return r.AddMethod("PUT") }

// WithQuery is a chainable helper. It sets the Query field to true.
func (r *RouteConfig) WithQuery() *RouteConfig {
	r.Query = true
	return r
}

// WithQuery is a chainable helper. It sets the Form field to true.
func (r *RouteConfig) WithForm() *RouteConfig {
	r.Form = true
	return r
}

// WithUser is a chainable helper. It sets the User field to true.
func (r *RouteConfig) WithUser() *RouteConfig {
	r.User = true
	return r
}

// WithBody is a chainable helper. It sets the Body field to true.
func (r *RouteConfig) WithBody() *RouteConfig {
	r.Body = true
	return r
}

// WithPrefix is a chainable helper. It sets the PathPrefix field to true.
func (r *RouteConfig) WithPrefix() *RouteConfig {
	r.PathPrefix = true
	return r
}

// RequireGroup is a chainable helper. It calls RequireGroup on the embedded
// Require.
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
		r.ID = r.String()
	}
	if strings.Contains(r.Path, "{") || strings.Contains(r.Host, "{") {
		r.PathVars = true
	}
	return nil
}

// String fulfills Stringer. Returns the RouteConfig as "(Method)path...", where
// the ellipse will only be present if the path is a prefix.
func (r *RouteConfig) String() string {
	prefxStr := ""
	if r.PathPrefix {
		prefxStr = "..."
	}
	host := ""
	if r.Host != "" {
		host = fmt.Sprintf("(HOST:%s)", r.Host)
	}
	return fmt.Sprintf("(%s) %s%s%s", r.Method, host, r.Path, prefxStr)
}

func (r *RouteConfig) Methods() []string {
	out := strings.Split(r.Method, ",")
	for i, m := range out {
		out[i] = strings.TrimSpace(m)
	}
	return out
}
