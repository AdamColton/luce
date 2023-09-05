package lstr

// CamelCase logic is used to split a string into words.
func CamelCase(str string) SubStrings {
	return TransitionSplit(IsUpper, str)
}

// TransitionSplit takes a Matcher and each place where the result transitions
// from false to true there is a split.
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
