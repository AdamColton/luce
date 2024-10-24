package handler_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/util/handler"
	"github.com/stretchr/testify/assert"
)

type point struct{ X, Y int }

func TestDeserialize(t *testing.T) {
	cmds, err := handler.Cmds([]*handler.Command{
		{Name: "move", Action: func(p *point) string { return fmt.Sprintf("%d,%d", p.X, p.Y) }},
		{Name: "sum", Action: func(p point) int { return p.X + p.Y }},
		{Name: "double", Action: func(i int) int { return i * 2 }},
		{Name: "ping", Action: func() string { return "pong" }},
		{Name: "fail", Action: func(i int) (int, error) { return 0, lerr.Str("failed") }},
		{
			Name:   "group",
			Action: func() {},
			Subcmds: []*handler.Command{
				{Name: "neg", Action: func(i int) int { return -i }},
			},
		},
	})
	assert.NoError(t, err)

	tt := map[string]struct {
		path     []string
		data     string
		expected any
		err      string
	}{
		"pointer":       {[]string{"move"}, `{"X":1,"Y":2}`, "1,2", ""},
		"struct":        {[]string{"sum"}, `{"X":3,"Y":4}`, 7, ""},
		"int":           {[]string{"double"}, `21`, 42, ""},
		"no argument":   {[]string{"ping"}, `1`, nil, "command takes no argument: ping"},
		"below a group": {[]string{"group", "neg"}, `5`, -5, ""},
		"handler error": {[]string{"fail"}, `1`, 0, "failed"},
		"bad data":      {[]string{"move"}, `not json`, nil, "invalid character 'o' in literal null (expecting 'u')"},
		"unknown":       {[]string{"nope", "x"}, `1`, nil, "Could not find command: nope/x"},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := cmds.Deserialize(json.Deserializer{}, []byte(tc.data), tc.path...)
			assert.Equal(t, tc.expected, got)
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err)
			}
		})
	}
}
