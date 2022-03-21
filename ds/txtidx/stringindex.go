package txtidx

type StringIndex struct {
	Strings []*String
	*Source
	Lookup map[uint32][]uint32 //rootIndex to StringsIndexes
}

func NewStringIndex() *StringIndex {
	return &StringIndex{
		Source: NewSource(),
		Lookup: make(map[uint32][]uint32),
	}
}

func (si *StringIndex) Get(idx int) string {
	s := si.Strings[idx]
	return Build(si.Source, s)
}

func (si *StringIndex) AddString(s *String) {
	idx := uint32(len(si.Strings))
	for _, iw := range s.IndexWords {
		si.Lookup[iw.BaseID] = append(si.Lookup[iw.BaseID], idx)
	}
}
