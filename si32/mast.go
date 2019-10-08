package si32

import (
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

// Build constructs an FST virtual machine from the given inputs.
func Build(src PairSlice) (t FST, err error) {
	m := BuildMAST(src)
	return m.BuildFST()
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

// BuildMAST builds a minimal acyclic subsequential transducer from the given inputs.
func BuildMAST(input PairSlice) (m MAST) {
	if !sort.IsSorted(input) {
		sort.Sort(input)
	}

	const initialMASTSize = 1024
	dic := make(map[int64][]*State)
	m.States = make([]*State, 0, initialMASTSize)
	m.FinalStates = make([]*State, 0, initialMASTSize)

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
				m.AddState(s)
				dic[s.hcode] = append(dic[s.hcode], s)
			}
			buf[i].Clear()
			buf[i-1].SetTransition(prev[i-1], s)
		}
		for i, size := prefixLen+1, len(in); i <= size; i++ {
			buf[i-1].SetTransition(in[i-1], buf[i])
		}
		if in != prev {
			buf[len(in)].IsFinal = true
		}
		for j := 1; j < prefixLen+1; j++ {
			outSuff, ok := buf[j-1].Output[in[j-1]]
			if ok {
				if outSuff == out {
					out = 0
					break
				}
				buf[j-1].RemoveOutput(in[j-1]) // clear the prev edge
				for ch := range buf[j].Trans {
					buf[j].SetOutput(ch, outSuff)
				}
				if buf[j].IsFinal && outSuff != 0 {
					buf[j].AddTail(outSuff)
				}
			}
		}
		if in != prev {
			buf[prefixLen].SetOutput(in[prefixLen], out)
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
			m.AddState(s)
			dic[s.hcode] = append(dic[s.hcode], s)
		}
		buf[i-1].SetTransition(prev[i-1], s)
	}
	m.StartingState = buf[0]
	m.AddState(buf[0])

	return
}

// Run rus the transducer in the given input.
func (m *MAST) Run(input string) (out []int32, ok bool) {
	s := m.StartingState
	for i, size := 0, len(input); i < size; i++ {
		if o, ok := s.Output[input[i]]; ok {
			out = append(out, o)
		}
		if s, ok = s.Trans[input[i]]; !ok {
			return
		}
	}
	for _, t := range s.Tails() {
		out = append(out, t)
	}
	return
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

// Dot outputs the FST in graphviz format.
func (m *MAST) Dot(w io.Writer) {
	fmt.Fprintln(w, "digraph G {")
	fmt.Fprintln(w, "\trankdir=LR;")
	fmt.Fprintln(w, "\tnode [shape=circle]")
	for _, s := range m.FinalStates {
		fmt.Fprintf(w, "\t%d [peripheries = 2];\n", s.ID)
	}
	for _, from := range m.States {
		for in, to := range from.Trans {
			fmt.Fprintf(w, "\t%d -> %d [label=\"%02X/%v", from.ID, to.ID, in, from.Output[in])
			if to.HasTail() {
				fmt.Fprintf(w, " %v", to.Tails())
			}
			fmt.Fprintln(w, "\"];")
		}
	}
	fmt.Fprintln(w, "}")
}
