package type32_test

type person struct {
	Name string
	Age  int
}

func (*person) TypeID32() uint32 {
	return 12345
}
