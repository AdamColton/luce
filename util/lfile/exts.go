package lfile

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// FilterHidden is a regular expression that matches a path if its last part
	// begins with a dot, so it can be used as the skipDir of a Match to skip
	// hidden directories.
	FilterHidden = `\/\.\.?[^\/]+$`
)

// Exts builds a regular expression that matches the paths that end in one of
// the extensions, given without the dot. The extensions are matched as they are
// written, so "tar.gz" only matches ".tar.gz". If hiddenDirs is false, paths
// that have a parent directory that is hidden (its name begins with a dot) are
// left out, though "." and ".." are not hidden.
func Exts(hiddenDirs bool, exts ...string) string {
	pre := ".*"
	if !hiddenDirs {
		//pre = `([^\/]|(\/[^\.]))*`
		pre = `^\/?(((\.\.?)|([^\/\.][^\/]*))\/)*[^\/]*`
	}
	quoted := make([]string, len(exts))
	for i, ext := range exts {
		quoted[i] = regexp.QuoteMeta(ext)
	}
	return fmt.Sprintf(`%s\.((%s))$`, pre, strings.Join(quoted, ")|("))
}
