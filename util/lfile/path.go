package lfile

// Name returns the last portion of a path as it's name.
//   - "/foo/bar.txt" => "/foo/","bar.txt"
//   - "/foo/bar/" => "/foo/","bar"
//   - "foo.txt" => "", "foo.txt"
//   - foo/ => "", "foo"
//
// The second returned value is the name and the first
// is the preceeding portion.
func Name(path string) (string, string) {
	end := len(path) - 1
	if end < 0 {
		return "", ""
	}
	for end > 0 && path[end] == '/' {
		end--
	}
	start := end - 1
	for ; start >= 0 && path[start] != '/'; start-- {
	}
	start++
	return path[:start], path[start : end+1]
}
