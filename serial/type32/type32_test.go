package type32_test

type person struct {
	Name string
	Age  int
}

func (*person) TypeID32() uint32 {
	return 12345
}

type strSlice []string

func (strSlice) TypeID32() uint32 {
	return 67890
}
