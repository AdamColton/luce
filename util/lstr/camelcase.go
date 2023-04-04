package lstr

// CamelCase logic is used to split a string into words at the upper case
// letters. For instance "helloWorldFoo" splits into "hello", "World" and
// "Foo". An upper case run at the start or the end of the string is not split
// ("HTTPServer" stays whole). Inside the string, a run of two or more upper case
// letters is split so its last letter starts the next word ("fooBARbaz" splits
// into "foo", "BA" and "Rbaz"). The indexes are byte indexes, so a multi-byte
// rune next to the split can be cut in two.
func CamelCase(str string) SubStrings {
	return TransitionSplit(IsUpper, str)
}

// TransitionSplit takes a Matcher and each place where the result transitions
// from false to true there is a split. CamelCase is TransitionSplit with the
// IsUpper Matcher, see it for the details.
func TransitionSplit(m Matcher, str string) SubStrings {
	lastL := -1
	var splits []int
	for s := NewScanner(str); !s.Done(); s.Next() {
		isU := s.Peek(m)
		if !isU {
			if lastL+2 <= s.I && lastL > -1 {
				if lastL+2 < s.I {
					splits = append(splits, lastL+1)
				}
				splits = append(splits, s.I-1)
			}
			lastL = s.I
		}
	}
	splits = append(splits, len(str))
	return SubStringBySplit(splits)
}
