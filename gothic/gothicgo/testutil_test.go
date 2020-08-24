package gothicgo

func newFile(names ...string) (*MemoryContext, *File) {
	pkg, file := names[0], names[0]
	if len(names) > 1 {
		file = names[1]
	}
	ctx := NewMemoryContext()
	return ctx, ctx.MustPackage(pkg).File(file)
}
