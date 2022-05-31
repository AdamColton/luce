package bintrie

type Multi[U Uint] []Trie[U]

func (m Multi[U]) Has(u U) bool {
	ln := len(m)
	mn := make([]*node[U], ln)
	for i, t := range m {
		mn[i] = t.(*node[U])
	}

	rm := func(idx int) {
		ln--
		if idx < ln {
			mn[idx] = mn[ln]
		}
		mn = mn[:ln]
	}

	prune := func(b U) {
		for j := 0; j < ln; {
			if mn[j].branches[b] == nil {
				rm(j)
			} else {
				mn[j] = mn[j].branches[b]
				j++
			}
		}
	}

	for i := U(0); i < sizeOf(u); i++ {
		b := u & 1
		u >>= 1
		prune(b)
		if ln == 0 {
			return false
		}
	}
	for _, n := range mn {
		if n.terminal {
			return true
		}
	}
	return false
}
