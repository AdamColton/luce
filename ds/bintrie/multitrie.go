package bintrie

type Multi []Trie

func (m Multi) Has(u uint32) bool {
	ln := len(m)
	mn := make([]*node, ln)
	for i, t := range m {
		mn[i] = t.(*node)
	}

	rm := func(idx int) {
		ln--
		if idx < ln {
			mn[idx] = mn[ln]
		}
		mn = mn[:ln]
	}

	prune := func(b uint32) {
		for j := 0; j < ln; {
			if mn[j].branches[b] == nil {
				rm(j)
			} else {
				mn[j] = mn[j].branches[b]
				j++
			}
		}
	}

	for i := 0; i < 32; i++ {
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
