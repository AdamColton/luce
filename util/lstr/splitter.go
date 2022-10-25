package lstr

type Splitter interface {
	Split(string) SubStrings
}

type SplitterFunc func(s string) SubStrings

func (fn SplitterFunc) Split(s string) SubStrings {
	return fn(s)
}
