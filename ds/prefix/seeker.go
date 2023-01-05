package prefix

type seeker struct {
	runes []rune
	idx   int
	p     *Prefix
	n     *node
}

func (s *seeker) moveNext(insert bool) (done bool) {
	done = s.n == nil || s.idx >= len(s.runes)
	if !done {
		r := s.runes[s.idx]
		s.idx++
		n := s.n.children.GetVal(r)
		if n == nil {
			done = !insert
			if insert {
				n = newNode()
				n.r = r
				n.parent = s.n
				s.n.children.Set(r, n)
				s.p.starts[r] = append(s.p.starts[r], n)
			}
		}
		s.n = n
	}
	return
}
