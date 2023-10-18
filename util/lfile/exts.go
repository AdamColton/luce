package lfile

import (
	"fmt"
	"strings"
)

const (
	FilterHidden = `\/\.[^\/]*$`
)

// Exts builds a regular expression for file extensions. If hiddenDirs is false
// files that have a parent directory that is hidden (name begins with .) those
// will be ommited.
func Exts(hiddenDirs bool, exts ...string) string {
	pre := ".*"
	if !hiddenDirs {
		pre = `([^\/]|(\/[^\.]))*`
	}
	return fmt.Sprintf(`^%s\.((%s))$`, pre, strings.Join(exts, ")|("))
}
