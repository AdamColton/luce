package gothicgo

// FuncRef represents a function as a type.
type FuncRef struct {
	*FuncType
	Pkg PackageRef
}

// NewFuncRef creates a FuncRef representing a Func as a Type.
func NewFuncRef(pkg PackageRef, name string, args ...NameType) *FuncRef {
	return &FuncRef{
		FuncType: NewFuncType(name, args...),
		Pkg:      pkg,
	}
}

// Call produces a invocation of the function and fulfills the FuncCaller
// interface
func (f *FuncRef) Call(pre Prefixer, args ...string) string {
	return funcCall(pre, f.Name, args, f.Pkg)
}
