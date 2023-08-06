package prefix

import "github.com/adamcolton/luce/util/navigator"

type seeker struct {
	p *Prefix
	*navigator.Navigator[rune, *node, *Prefix]
}
