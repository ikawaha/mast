package string2integer

type dict struct {
	bucket map[uint][]*state
}

func newDict() *dict {
	return &dict{
		bucket: map[uint][]*state{},
	}
}

func (d dict) find(s *state) (*state, bool) {
	states, ok := d.bucket[s.hashCode]
	if !ok {
		return nil, false
	}
	for _, v := range states {
		if v.equal(s) {
			return v, true
		}
	}
	return nil, false
}

func (d *dict) add(s *state) {
	d.bucket[s.hashCode] = append(d.bucket[s.hashCode], s)
}
