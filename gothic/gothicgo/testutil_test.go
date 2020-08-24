package gothicgo

func newFile(name string) (*MemoryContext, *File) {
	ctx := NewMemoryContext()
	return ctx, ctx.MustPackage(name).File("foo")
}
