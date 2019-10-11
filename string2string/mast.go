package string2string

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
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

// BuildMAST builds a minimal acyclic subsequential transducer from the given inputs.
func BuildMAST(input PairSlice) *MAST {
	sort.Sort(input)

	const initialMASTSize = 1024
	dic := make(map[uint][]*State)
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
			s.SetInvTransition()
		}
		for i := prefixLen + 1; i <= len(in); i++ {
			buf[i-1].SetTransition(in[i-1], buf[i])
		}
		if in != prev {
			buf[len(in)].IsFinal = true
		}
		for j := 1; j < prefixLen+1; j++ {
			outPrefix := commonPrefix(buf[j-1].Output[in[j-1]], out)
			outSuffix := strings.TrimPrefix(buf[j-1].Output[in[j-1]], outPrefix)
			buf[j-1].SetOutput(in[j-1], outPrefix)
			for ch := range buf[j].Trans {
				buf[j].SetOutput(ch, outSuffix+buf[j].Output[ch])
			}
			if buf[j].IsFinal {
				set := stringSet{}
				if !buf[j].HasTail() {
					set[outSuffix] = struct{}{}
				} else {
					for _, s := range buf[j].Tails() {
						s = outSuffix + s
						set[s] = struct{}{}
					}
				}
				buf[j].Tail = set
			}
			out = strings.TrimPrefix(out, outPrefix)
		}
		if in == prev {
			buf[len(in)].AddTail(out)
		} else {
			buf[prefixLen].SetOutput(in[prefixLen], out)
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
		}
		buf[i-1].SetTransition(prev[i-1], s)
		s.SetInvTransition()
	}
	ret.StartingState = buf[0]
	ret.AddState(buf[0])

	return &ret
}

// Run rus the transducer in the given input.
func (m *MAST) Run(input string) (out []string, accept bool) {
	var buf bytes.Buffer
	s := m.StartingState
	for i, size := 0, len(input); i < size; i++ {
		if o, ok := s.Output[input[i]]; ok {
			buf.WriteString(o)
		}
		var ok bool
		s, ok = s.Trans[input[i]]
		if !ok {
			return out, false
		}
	}
	o := buf.String()
	if !s.HasTail() {
		out = append(out, o)
		return out, true
	}
	for _, t := range s.Tails() {
		out = append(out, o+t)
	}
	return out, true
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
