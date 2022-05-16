package lhttp

// Join elems making sure there is a single / between each elem.
func Join(elems ...string) string {
	if len(elems) == 0 {
		return ""
	}
	if len(elems) == 1 {
		return elems[0]
	}

	prev := elems[0]
	pln := len(prev)
	ln := pln
	for _, e := range elems[1:] {
		if e == "" {
			continue
		}
		ps := prev[pln-1] == '/'
		es := e[0] == '/'
		prev, pln = e, len(e)
		ln += pln
		if ps && es {
			ln--
		} else if !ps && !es {
			ln++
		}
	}

	out := make([]byte, 0, ln)

	out = append(out, elems[0]...)
	for _, e := range elems[1:] {
		if e == "" {
			continue
		}
		ps := out[len(out)-1] == '/'
		es := e[0] == '/'
		if ps && es {
			out = append(out, e[1:]...)
		} else {
			if !ps && !es {
				out = append(out, '/')
			}
			out = append(out, e...)
		}
	}
	return string(out)
}
