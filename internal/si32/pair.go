package si32

// Pair implements a pair of input and output.
type Pair struct {
	In  string
	Out int32
}

// PairSlice implements a slice of input and output pairs.
type PairSlice []Pair

func (ps PairSlice) Len() int      { return len(ps) }
func (ps PairSlice) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }
func (ps PairSlice) Less(i, j int) bool {
	if ps[i].In == ps[j].In {
		return ps[i].Out < ps[j].Out
	}
	return ps[i].In < ps[j].In
}

func (ps PairSlice) maxInputWordLen() (max int) {
	for _, pair := range ps {
		if size := len(pair.In); size > max {
			max = size
		}
	}
	return
}
