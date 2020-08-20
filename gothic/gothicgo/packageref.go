package gothicgo

import (
	"regexp"
	"strings"

	"github.com/adamcolton/luce/lerr"
)

// PackageRef represents a reference to a package.
type PackageRef interface {
	ImportPath() string
	Name() string

	// ImportSpec is the specification and may include a package name or a
	// modifier like _.
	ImportSpec() string

	// PackageRef is not meant to be implemented, it's meant as an accessor to the
	// underlying packageRef. All instances should be created with NewPackageRef
	// to guarentee that the reference is well formed.
	privatePkgRef()
}

var pkgBuiltin = &packageRef{}

// PkgBuiltin returns a PackageRef with an empty name.
func PkgBuiltin() PackageRef {
	return pkgBuiltin
}

type packageRef struct {
	path, name string
}

func (p *packageRef) ImportPath() string {
	return p.path
}

func (p *packageRef) Name() string {
	return p.name
}

func (p *packageRef) ImportSpec() string {
	return strings.Join([]string{"\"", p.path, "\""}, "")
}

func (*packageRef) privatePkgRef() {}

// TODO: this regex is only mostly right
var packageRefRegex = regexp.MustCompile(`^(?:[\w\-\.]+\/)*([\w\-]+)$`)

// ErrBadPackageRef indicates a poorly formatted package ref string.
const ErrBadPackageRef = lerr.Str("Bad Package Ref")

// NewPackageRef takes the string used to import a pacakge and returns a
// PackageRef.
func NewPackageRef(ref string) (PackageRef, error) {
	m := packageRefRegex.FindStringSubmatch(ref)
	if len(m) == 0 {
		return nil, ErrBadPackageRef
	}
	return &packageRef{
		path: m[0],
		name: m[1],
	}, nil
}

// MustPackageRef returns a new PackageRef and panics if there is an error
func MustPackageRef(ref string) PackageRef {
	p, err := NewPackageRef(ref)
	lerr.Panic(err)
	return p
}
