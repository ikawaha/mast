package string2integer

// Pair implements a pair of input and output.
type Pair struct {
	In  string
	Out int
}

// PairList implements a slice of input and output pairs.
type PairList []Pair

func (ps PairList) Len() int {
	return len(ps)
}

func (ps PairList) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps PairList) Less(i, j int) bool {
	if ps[i].In == ps[j].In {
		return ps[i].Out < ps[i].Out
	}
	return ps[i].In < ps[j].In
}

func (ps PairList) maxInputWordLen() (max int) {
	for _, pair := range ps {
		if size := len(pair.In); size > max {
			max = size
		}
	}
	return
}
