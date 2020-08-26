package main

import (
	"text/template"

	"github.com/adamcolton/luce/util/luceio"

	"github.com/adamcolton/luce/gothic/gothicgo"
	"github.com/adamcolton/luce/gothic/gothicgo/ggutil"
)

var (
	test = template.Must(template.New("test").Parse(`func Test{{.Name}}TypeGen(t *testing.T){
		x := {{.Constructor}}

		assert.Equal(t, {{.Kind}}, x.Kind())
		assert.Equal(t, {{.Package}}, x.PackageRef())

		n := x.Named("Foo")
		assert.Equal(t, "Foo", n.Name())
		assert.Equal(t, x, n.Type)
	
		n = x.Unnamed()
		assert.Equal(t, "", n.Name())
		assert.Equal(t, x, n.Type)
	
		p := x.Ptr()
		assert.Equal(t, PointerKind, p.Kind())
		assert.Equal(t, x, p.Elem())
	
		s := x.Slice()
		assert.Equal(t, SliceKind, s.Kind())
		assert.Equal(t, x, s.Elem())
	
		a := x.Array(13)
		assert.Equal(t, ArrayKind, a.Kind())
		assert.Equal(t, x, a.Elem())
	
		mp := x.AsMapElem(IntType)
		assert.Equal(t, MapKind, mp.Kind())
		assert.Equal(t, x, mp.Elem())
	
		mp = x.AsMapKey(IntType)
		assert.Equal(t, MapKind, mp.Kind())
		assert.Equal(t, x, mp.MapKey())

		str := PrefixWriteToString(x, DefaultPrefixer)
		assert.Equal(t, "{{.String}}", str)
	}`))

	typegen = template.Must(template.New("typegen").Parse(`
	// Named fulfills Type. Returns a NameType with the given name.
	func ({{.R}} {{.Name}}) Named(name string) NameType { return NameType{name, {{.R}}} }

	// Unnamed funfills Type. Returns a NameType with an empty Name.
	func ({{.R}} {{.Name}}) Unnamed() NameType { return NameType{"", {{.R}}} }

	// Ptr funfills Type.
	func ({{.R}} {{.Name}}) Ptr() PointerType { return PointerTo({{.R}}) }

	// Slice funfills Type.
	func ({{.R}} {{.Name}}) Slice() SliceType { return SliceOf({{.R}}) }

	// Array funfills Type.
	func ({{.R}} {{.Name}}) Array(size int) ArrayType { return ArrayOf({{.R}}, size) }

	// AsMapElem funfills Type.
	func ({{.R}} {{.Name}}) AsMapElem(key Type) MapType { return MapOf(key, {{.R}}) }

	// AsMapKey funfills Type. Returns a NameType with an empty Name.
	func ({{.R}} {{.Name}}) AsMapKey(elem Type) MapType { return MapOf({{.R}}, elem) }
`))
)

type fullType struct {
	R, Name string

	// Test values
	Constructor, String, Kind, Package string
}

func main() {
	pth := "github.com/adamcolton/luce/gothic/"
	ctx := gothicgo.ContextFactory{
		OutputPath:     ggutil.GoSrc(pth),
		ImportPath:     pth,
		DefaultComment: "Generated by gothicgo/bootstrap - Do not modify",
	}.New()
	pkg := ctx.MustPackage("gothicgo")
	typeFile := pkg.File("type.gen")
	testFile := pkg.File("type.gen_test")

	testFile.Imports.Add(
		gothicgo.MustExternalPackageRef("testing"),
		gothicgo.MustExternalPackageRef("github.com/testify/assert"),
	)

	fts := []fullType{
		{"a", "ArrayType", "IntType.Array(5)", "[5]int", "ArrayKind", "PkgBuiltin()"},
		{"b", "builtin", "IntType", "int", "IntKind", "PkgBuiltin()"},
	}
	for _, ft := range fts {
		typeFile.AddWriterTo(&luceio.TemplateTo{
			TemplateExecutor: typegen,
			Data:             ft,
		})
		testFile.AddWriterTo(&luceio.TemplateTo{
			TemplateExecutor: test,
			Data:             ft,
		})
	}

	ctx.MustExport()
}
