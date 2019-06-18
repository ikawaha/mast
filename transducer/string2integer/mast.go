package string2integer

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

const (
	initialMastSize = 1024
)

// Mast represents Minimal Acyclic Subsequential Transducers.
type Mast struct {
	initialState *state
	states       []*state
	finalStates  []*state
}

func newMast() *Mast {
	return &Mast{
		states:      make([]*state, 0, initialMastSize),
		finalStates: make([]*state, 0, initialMastSize),
	}
}

func (m *Mast) addState(n *state) {
	n.ID = len(m.states)
	m.states = append(m.states, n)
	if n.isFinal {
		m.finalStates = append(m.finalStates, n)
	}
}

func commonPrefix(a, b string) string {
	end := len(a)
	if end > len(b) {
		end = len(b)
	}
	i := 0
	for i < end && a[i] == b[i] {
		i++
	}
	return a[0:i]
}

func BuildMast(inputs PairList) (m *Mast) {
	m = newMast()
	dict := newDict()
	buf := make([]*state, inputs.maxInputWordLen()+1)
	for i := range buf {
		buf[i] = newState()
	}

	sort.Sort(inputs)

	prev := ""
	for _, v := range inputs {
		in, out := v.In, v.Out
		prefixLen := len(commonPrefix(in, prev))
		for i := len(prev); i > prefixLen; i-- {
			s, ok := dict.find(buf[i])
			if !ok {
				s = &state{}
				*s = *buf[i] // shallow copy
				m.addState(s)
				dict.add(s)
			}
			buf[i].clear()
			buf[i-1].setTransition(prev[i-1], s)
			s.setInvTransition()
		}
		for i, size := prefixLen+1, len(in); i <= size; i++ {
			buf[i-1].setTransition(in[i-1], buf[i])
		}
		if in != prev {
			buf[len(in)].isFinal = true
		}
		for i := 1; i < prefixLen+1; i++ {
			ch := in[i-1]
			if v, ok := buf[i-1].outputs[ch]; ok && v != out {
				buf[i-1].deleteOutput(ch)
				for ch := range buf[i].transition {
					buf[i].setOutput(ch, v)
				}
				if buf[i].isFinal {
					buf[i].addTail(v)
				}
			}
		}
		if in == prev {
			buf[len(in)].addTail(out)
		} else {
			buf[prefixLen].setOutput(in[prefixLen], out)
		}
		prev = in
	}
	// flush the buf
	for i := len(prev); i > 0; i-- {
		s, ok := dict.find(buf[i])
		if !ok {
			s = &state{}
			*s = *buf[i] // shallow copy
			m.addState(s)
			dict.add(s)
		}
		buf[i-1].setTransition(prev[i-1], s)
		s.setInvTransition()
	}
	m.initialState = buf[0]
	m.addState(buf[0])

	return
}

func (m *Mast) Run(input string) ([]int, bool) {
	ret := []int{}
	s := m.initialState
	for i, size := 0, len(input); i < size; i++ {
		if v, ok := s.outputs[input[i]]; ok {
			ret = append(ret, v)
		}
		next, ok := s.transition[input[i]]
		if !ok {
			return nil, false
		}
		s = next
	}
	for _, v := range s.tails() {
		ret = append(ret, v)
	}
	return ret, true
}

func (m *Mast) Accept(input string) bool {
	s := m.initialState
	for i, size := 0, len(input); i < size; i++ {
		next, ok := s.transition[input[i]]
		if !ok {
			return false
		}
		s = next
	}
	return true
}

func (m *Mast) Dot(w io.Writer) error {
	if _, err := fmt.Fprintln(w, "digraph G {"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "\trankdir=LR;"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "\tnode [shape=circle]"); err != nil {
		return err
	}
	for _, s := range m.finalStates {
		if _, err := fmt.Fprintf(w, "\t%d [peripheries = 2];\n", s.ID); err != nil {
			return err
		}
	}
	for _, from := range m.states {
		for in, to := range from.transition {
			var out string
			if v, ok := from.outputs[in]; ok {
				out = strconv.Itoa(v)
			}
			if _, err := fmt.Fprintf(w, "\t%d -> %d [label=\"%c/%v", from.ID, to.ID, in, out); err != nil {
				return err
			}
			if to.hasTail() {
				if _, err := fmt.Fprintf(w, " %v", to.tails()); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w, "\"];"); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintln(w, "}"); err != nil {
		return err
	}
	return nil
}
