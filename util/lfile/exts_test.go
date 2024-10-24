package lfile

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExts(t *testing.T) {
	ext := Exts(false, "js")
	re := regexp.MustCompile(ext)
	assert.True(t, re.MatchString("../../foo/test.js"))
	assert.False(t, re.MatchString("../.hidden/foo/test.js"))
	assert.False(t, re.MatchString(".hidden/nothidden/foo/test.js"))

	tt := map[string]struct {
		exts string
		pass []string
		fail []string
	}{
		"js": {
			exts: Exts(false, "js"),
			pass: []string{
				"foo.js",
				"foo/bar.js",
				"/a/b/cde.js",
				"./hello.js",
				"../hello/goodbye.js",
			},
			fail: []string{
				"foo.txt",
				"foo.js.txt",
				".hidden/foo.js",
				"/.hidden/foo.js",
				"/nothidden/.hidden/foo.js",
			},
		},
		"js-withHidden": {
			exts: Exts(true, "js"),
			pass: []string{
				"foo.js",
				"foo/bar.js",
				"/a/b/cde.js",
				"./hello.js",
				"../hello/goodbye.js",
				"/.hidden/foo.js",
				"/nothidden/.hidden/foo.js",
			},
			fail: []string{
				"foo.txt",
				"foo.js.txt",
			},
		},
		"js&txt": {
			exts: Exts(false, "js", "txt"),
			pass: []string{
				"foo.js",
				"foo/bar.js",
				"/a/b/cde.js",
				"./hello.js",
				"../hello/goodbye.js",
				"foo.txt",
				"foo.js.txt",
			},
			fail: []string{
				"/.hidden/foo.js",
				"/nothidden/.hidden/foo.js",
				"foo.js.txt.other",
				"foo.css",
			},
		},
	}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			re := regexp.MustCompile(tc.exts)
			for _, s := range tc.pass {
				t.Run(s, func(t *testing.T) {
					assert.True(t, re.MatchString(s))
				})
			}
			for _, s := range tc.fail {
				t.Run(s, func(t *testing.T) {
					assert.False(t, re.MatchString(s))
				})
			}
		})
	}
}
