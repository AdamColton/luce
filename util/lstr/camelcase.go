package lstr

const (
	lower byte = iota + 1
	upper
	digit
)

// AbcdEfgh -> Abcd Efgh
// AbcdEFgh -> Abcd E Fgh
// AbcDEFgh -> Abc DE Fgh

// P | C |
// U | U | AB
// U | L | Ab
// L | U | aB
// L | L | ab

func CamelCase(str string) SubStrings {
	// lastU := -1
	// lastL := -1
	prev := false
	prevIdx := -1
	var splits []int
	for s := NewScanner(str); !s.Done(); s.Next() {
		cur := s.Peek(IsUpper)
		// if cur {
		// 	lastU = s.Idx
		// } else {
		// 	lastL = s.Idx
		// }
		if prev && !cur && prevIdx > 0 {
			splits = append(splits, prevIdx)
		}
		prev = cur
		prevIdx = s.Idx
	}
	splits = append(splits, len(str))
	// U -> U : scan until l or EOF, split one back "HTMLLoader"
	// U -> l : scan until U or EOF, "HtmlLoader"
	//
	return SubStringBySplit(splits)
}
