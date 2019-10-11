package string2int32

import (
	"bufio"
	"fmt"
	"io"
	"sort"
)

// MAST represents a Minimal Acyclic Subsequential Transducer.
type MAST struct {
	StartingState *State
	States        []*State
	FinalStates   []*State
}

// AddState adds a state to the finite state transducer.
func (m *MAST) AddState(n *State) {
	n.ID = len(m.States)
	m.States = append(m.States, n)
	if n.IsFinal {
		m.FinalStates = append(m.FinalStates, n)
	}
}

func commonPrefix(a, b string) string {
	end := len(a)
	if end > len(b) {
		end = len(b)
	}
	var i int
	for i < end && a[i] == b[i] {
		i++
	}
	return a[0:i]
}

const initialMASTSize = 1024

// BuildMAST builds a minimal acyclic subsequential transducer from the given inputs.
func BuildMAST(input PairSlice) *MAST {
	sort.Sort(input)

	dic := make(map[int64][]*State)
	ret := MAST{
		States:      make([]*State, 0, initialMASTSize),
		FinalStates: make([]*State, 0, initialMASTSize),
	}

	buf := make([]*State, input.maxInputWordLen()+1)
	for i := range buf {
		buf[i] = NewState()
	}
	prev := ""
	for _, pair := range input {
		in, out := pair.In, pair.Out
		fZero := (out == 0) // flag
		prefixLen := len(commonPrefix(in, prev))
		for i := len(prev); i > prefixLen; i-- {
			var s *State
			if cs, ok := dic[buf[i].hcode]; ok {
				for _, c := range cs {
					if c.Equal(buf[i]) {
						s = c
						break
					}
				}
			}
			if s == nil {
				s = &State{}
				*s = *buf[i]
				ret.AddState(s)
				dic[s.hcode] = append(dic[s.hcode], s)
			}
			buf[i].Clear()
			buf[i-1].SetTransition(prev[i-1], s)
		}
		for i := prefixLen + 1; i <= len(in); i++ {
			buf[i-1].SetTransition(in[i-1], buf[i])
		}
		if in != prev {
			buf[len(in)].IsFinal = true
		}
		var outSuffix bool
		for j := 1; j < prefixLen+1; j++ {
			v, ok := buf[j-1].Output[in[j-1]]
			if !ok {
				continue
			}
			if v == out {
				outSuffix = true
				break
			}
			buf[j-1].RemoveOutput(in[j-1]) // clear the prev edge
			for ch := range buf[j].Trans {
				buf[j].SetOutput(ch, v)
			}
			if buf[j].IsFinal && v != 0 {
				buf[j].AddTail(v)
			}
		}
		if in != prev {
			if !outSuffix {
				buf[prefixLen].SetOutput(in[prefixLen], out)
			}
		} else if fZero || out != 0 {
			buf[len(in)].AddTail(out)
		}
		prev = in
	}
	// flush the buf
	for i := len(prev); i > 0; i-- {
		var s *State
		if cs, ok := dic[buf[i].hcode]; ok {
			for _, c := range cs {
				if c.Equal(buf[i]) {
					s = c
					break
				}
			}
		}
		if s == nil {
			s = &State{}
			*s = *buf[i]
			buf[i].Clear()
			ret.AddState(s)
			dic[s.hcode] = append(dic[s.hcode], s)
		}
		buf[i-1].SetTransition(prev[i-1], s)
	}
	ret.StartingState = buf[0]
	ret.AddState(buf[0])

	return &ret
}

// Run rus the transducer input the given input.
func (m *MAST) Run(input string) (out []int32, accept bool) {
	s := m.StartingState
	for i, size := 0, len(input); i < size; i++ {
		if o, ok := s.Output[input[i]]; ok {
			out = append(out, o)
		}
		var ok bool
		s, ok = s.Trans[input[i]]
		if !ok {
			return out, false
		}
	}
	return append(out, s.Tails()...), true
}

// Accept checks that the transducer accepts the given input.
func (m *MAST) Accept(input string) (ok bool) {
	s := m.StartingState
	for i, size := 0, len(input); i < size; i++ {
		if s, ok = s.Trans[input[i]]; !ok {
			return false
		}
	}
	return true
}

// Dot outputs the FST input graphviz format.
func (m *MAST) Dot(w io.Writer) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	bw.WriteString("digraph G {\n")
	bw.WriteString("\trankdir=LR;\n")
	bw.WriteString("\tnode [shape=circle]\n")
	for _, s := range m.FinalStates {
		fmt.Fprintf(bw, "\t%d [peripheries = 2];\n", s.ID)
	}
	for _, from := range m.States {
		for in, to := range from.Trans {
			if out, ok := from.Output[in]; !ok {
				fmt.Fprintf(bw, "\t%d -> %d [label=\"%02X(%c)/ε", from.ID, to.ID, in, in)
			} else {
				fmt.Fprintf(bw, "\t%d -> %d [label=\"%02X(%c)/%v", from.ID, to.ID, in, in, out)
			}
			if to.HasTail() {
				fmt.Fprintf(bw, " %v", to.Tails())
			}
			bw.WriteString("\"];\n")
		}
	}
	bw.WriteString("}\n")
}
