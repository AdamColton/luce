package ggutil

import (
	"go/build"
	"os"
	"path"
)

var gpVal string

func gp(elem ...string) []string {
	if gpVal == "" {
		gpVal = os.Getenv("GOPATH")
		if gpVal == "" {
			gpVal = build.Default.GOPATH
		}
	}
	return append([]string{gpVal}, elem...)
}

func GoPath(elem ...string) string {
	return path.Join(gp(elem...)...)
}

func GoSrc(elem ...string) string {
	return path.Join(append(gp("src"), elem...)...)
}

func Github(elem ...string) string {
	return path.Join(append(gp("src", "github.com"), elem...)...)
}
