package service

// func TestRouteConfigGenPath(t *testing.T) {
// 	tt := map[string]struct {
// 		expected   string
// 		base, path string
// 	}{
// 		"slash-on-base": {
// 			expected: "foo/bar",
// 			base:     "foo/",
// 			path:     "bar",
// 		},
// 		"slash-on-path": {
// 			expected: "foo/bar",
// 			base:     "foo",
// 			path:     "/bar",
// 		},
// 		"slash-on-both": {
// 			expected: "foo/bar",
// 			base:     "foo/",
// 			path:     "/bar",
// 		},
// 		"slash-on-neither": {
// 			expected: "foo/bar",
// 			base:     "foo",
// 			path:     "bar",
// 		},
// 	}

// 	for n, tc := range tt {
// 		t.Run(n, func(t *testing.T) {
// 			rc := (&RouteConfigGen{
// 				Base: tc.base,
// 			}).Path(tc.path)
// 			assert.Equal(t, tc.expected, rc.Path)
// 		})
// 	}
// }
