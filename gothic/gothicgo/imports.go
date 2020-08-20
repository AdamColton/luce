package gothicgo

import (
	"io"
	"sort"

	"github.com/adamcolton/luce/ds/bufpool"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// Imports is a tool for managing imports. Imports can be defined by path or
// package and either may include an alias. The ResolvePackages method must be
// called to resolve any packages to refs.
type Imports struct {
	self PackageRef
	refs map[string]PackageRef
}

// NewImports sets up an instance of Imports.
func NewImports(self PackageRef) *Imports {
	return &Imports{
		self: self,
		refs: make(map[string]PackageRef),
	}
}

// Prefix returns the name or alias of the package reference if it is
// different from Imports.self. The name will either be a blank string or will
// end with a period.
func (i *Imports) Prefix(ref PackageRef) string {
	if (i != nil && i.self != nil && ref.ImportPath() == i.self.ImportPath()) || ref.Name() == "" {
		return ""
	}
	return i.GetRefName(ref) + "."
}

// Add takes PackageRefs and adds them as imports without aliases.
func (i *Imports) Add(refs ...PackageRef) {
	for _, ref := range refs {
		if ref == nil {
			continue
		}
		rs := ref.ImportPath()
		if rs == "" || (i.self != nil && rs == i.self.ImportPath()) {
			continue
		}
		if _, exists := i.refs[rs]; !exists {
			i.refs[rs] = ref
		}
	}
}

// AddImports takes another instance of Imports and adds all it's imports. This
// runs the risk of clobbering aliases.
func (i *Imports) AddImports(imports *Imports) {
	for path, alias := range imports.refs {
		i.refs[path] = alias
	}
}

// RemoveRef removes a reference.
func (i *Imports) RemoveRef(ref PackageRef) {
	if i != nil && ref != nil {
		delete(i.refs, ref.ImportPath())
	}
}

// GetRefName takes a package ref and returns the name it will be referenced by
// in the Import context. If the package is aliased it will return the alias,
// otherwise it will return the package name. If there is an unresolved name
// matching the PackageRef, it will be treated as resolving to the ref.
func (i *Imports) GetRefName(ref PackageRef) string {
	if i == nil {
		return ref.Name()
	}

	if reg, ok := i.refs[ref.ImportPath()]; ok {
		if name := reg.Name(); name != "" {
			return name
		}
	}

	return ref.Name()
}

// String returns the imports as Go code.
func (i *Imports) String() string {
	buf := bufpool.Get()
	i.WriteTo(buf)
	return bufpool.PutStr(buf)
}

// WriteTo writes the Go code to a writer
func (i *Imports) WriteTo(w io.Writer) (int64, error) {
	ln := len(i.refs)
	if ln == 0 {
		return 0, nil
	}
	sum := luceio.NewSumWriter(w)
	sum.WriteString("import (")

	refs := make([]string, 0, len(i.refs))
	for path := range i.refs {
		refs = append(refs, path)
	}
	sort.Strings(refs)

	for _, path := range refs {
		sum.WriteString("\n\t")
		sum.WriteString(i.refs[path].ImportSpec())
	}
	sum.WriteString("\n)\n")
	sum.Err = lerr.Wrap(sum.Err, "While writing imports:")
	return sum.Rets()
}
