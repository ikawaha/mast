package ss

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// mast represents Minimal Acyclic Subsequential Transeducers.
type mast struct {
	initialState *state
	states       []*state
	finalStates  []*state
}

func newMast() (m *mast) {
	const initialMastSize = 1024
	m = new(mast)
	m.states = make([]*state, 0, initialMastSize)
	m.finalStates = make([]*state, 0, initialMastSize)
	return
}

func (m *mast) addState(n *state) {
	n.ID = len(m.states)
	m.states = append(m.states, n)
	if n.IsFinal {
		m.finalStates = append(m.finalStates, n)
	}
}

// Build returns a virtual machine of a finite state transducer.
func Build(input PairSlice) (vm FstVM, err error) {
	m := buildMast(input)
	return m.compile()
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

func commonPrefixLen(a, b string) int {
	end := len(a)
	if end > len(b) {
		end = len(b)
	}
	var i int
	for i < end && a[i] == b[i] {
		i++
	}
	return i
}

func buildMast(input PairSlice) *mast {
	sort.Sort(input)
	mast := newMast()
	buf := make([]*state, input.maxInputWordLen()+1)
	for i := range buf {
		buf[i] = newState()
	}
	prev := ""
	for _, pair := range input {
		in, out := pair.In, pair.Out
		prefixLen := commonPrefixLen(in, prev)
		candidate := mast.finalStates
		for i := len(prev); i > prefixLen; i-- {
			var s *state
			detected := false
			if candidate != nil {
				for _, c := range candidate {
					if c.eq(buf[i]) {
						buf[i].renew()
						s = c
						candidate = c.Prev
						detected = true
						break
					}
				}
			}
			if !detected {
				candidate = nil
				s = &state{}
				*s = *buf[i]
				buf[i].renew()
				mast.addState(s)
			}
			buf[i-1].setTransition(prev[i-1], s)
			s.setInvTransition()
		}
		for i, size := prefixLen+1, len(in); i <= size; i++ {
			buf[i-1].setTransition(in[i-1], buf[i])
		}
		if in != prev {
			buf[len(in)].IsFinal = true
		}
		for j := 1; j < prefixLen+1; j++ {
			outPref := commonPrefix(buf[j-1].Output[in[j-1]], out)
			outSuff := strings.TrimPrefix(buf[j-1].Output[in[j-1]], outPref)
			buf[j-1].setOutput(in[j-1], outPref)
			for ch := range buf[j].Trans {
				buf[j].setOutput(ch, outSuff+buf[j].Output[ch])
			}
			if buf[j].IsFinal {
				set := make(map[string]bool)
				if !buf[j].hasTail() {
					set[outSuff] = true
				} else {
					for _, s := range buf[j].tails() {
						s = outSuff + s
						set[s] = true
					}
				}
				buf[j].setTail(set)
			}
			out = strings.TrimPrefix(out, outPref)
		}
		if in == prev {
			buf[len(in)].addTail(out)
		} else {
			buf[prefixLen].setOutput(in[prefixLen], out)
		}
		prev = in
	}
	// flush the buf
	candidate := mast.finalStates
	for i := len(prev); i > 0; i-- {
		var s *state
		detected := false
		if candidate != nil {
			for _, c := range candidate {
				if c.eq(buf[i]) {
					s = c
					candidate = c.Prev
					detected = true
					break
				}
			}
		}
		if !detected {
			candidate = nil
			s = buf[i]
			mast.addState(s)
		}
		buf[i-1].setTransition(prev[i-1], s)
		s.setInvTransition()
	}
	mast.initialState = buf[0]
	mast.addState(buf[0])

	return mast
}

func (m *mast) run(input string) (out []string, ok bool) {
	var buf bytes.Buffer
	s := m.initialState
	for i, size := 0, len(input); i < size; i++ {
		if o, ok := s.Output[input[i]]; ok {
			buf.WriteString(o)
		}
		if s, ok = s.Trans[input[i]]; !ok {
			return
		}
	}
	o := buf.String()
	if !s.hasTail() {
		out = append(out, o)
		return
	}
	for _, t := range s.tails() {
		out = append(out, o+t)
	}
	return
}

func (m *mast) accept(input string) (ok bool) {
	s := m.initialState
	for i, size := 0, len(input); i < size; i++ {
		if s, ok = s.Trans[input[i]]; !ok {
			return
		}
	}
	return
}

func (m *mast) dot(w io.Writer) {
	fmt.Fprintln(w, "digraph G {")
	fmt.Fprintln(w, "\trankdir=LR;")
	fmt.Fprintln(w, "\tnode [shape=circle]")
	for _, s := range m.finalStates {
		fmt.Fprintf(w, "\t%d [peripheries = 2];\n", s.ID)
	}
	for _, from := range m.states {
		for in, to := range from.Trans {
			fmt.Fprintf(w, "\t%d -> %d [label=\"%X/%s", from.ID, to.ID, in, from.Output[in])
			if to.hasTail() {
				fmt.Fprintf(w, " %v", to.tails())
			}
			fmt.Fprintln(w, "\"];")
		}
	}
	fmt.Fprintln(w, "}")
}

type byteSlice []byte

func (p byteSlice) Len() int           { return len(p) }
func (p byteSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p byteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func toBytes(x int) []byte {
	var (
		b   byte
		buf []byte
	)
	for x >= 256 {
		b = byte(x % 256)
		x >>= 8
		buf = append(buf, b)
	}
	buf = append(buf, byte(x))
	return buf
}

func (m *mast) compile() (vm FstVM, err error) {
	var tape bytes.Buffer
	var edges []byte
	addrMap := make(map[int]int)
	for _, s := range m.states {
		edges = edges[:0]
		for ch := range s.Trans {
			edges = append(edges, ch)
		}
		if len(edges) > 0 {
			sort.Sort(byteSlice(edges))
		}
		for i, size := 0, len(edges); i < size; i++ {
			inp := edges[size-1-i]
			next := s.Trans[inp]
			out := s.Output[inp]
			addr, ok := addrMap[next.ID]
			if !ok && !next.IsFinal {
				err = fmt.Errorf("next addr is undefined: state(%v), input(%X)", s.ID, inp)
				return
			}
			var op instOp
			if len(out) > 0 {
				if i == 0 {
					op = instOutputBreak
				} else {
					op = instOutput
				}
			} else if i == 0 {
				op = instBreak
			} else {
				op = instMatch
			}
			inst := byte(op)
			jump := len(vm.prog) - addr
			if len(out) != 0 {
				adr := tape.Len()
				dst := toBytes(adr)
				sz := byte(len(dst))
				vm.prog = append(vm.prog, dst...)
				vm.prog = append(vm.prog, sz)

				tape.WriteString(out)
				tape.WriteByte(byte(0x00))
			}
			if jump > 1 {
				dst := toBytes(jump)
				inst |= byte(len(dst))
				vm.prog = append(vm.prog, dst...)
			}
			vm.prog = append(vm.prog, inp)
			vm.prog = append(vm.prog, inst)
		}
		if s.IsFinal {
			inst := byte(instAccept)
			if len(s.Tail) > 0 {
				dst1 := toBytes(tape.Len())
				inst |= byte(len(dst1))
				for t := range s.Tail {
					tape.WriteString(t)
					tape.WriteByte(byte(0x00))
				}
				dst2 := toBytes(tape.Len())
				vm.prog = append(vm.prog, dst2...)
				vm.prog = append(vm.prog, byte(len(dst2)))
				vm.prog = append(vm.prog, dst1...)
			}
			vm.prog = append(vm.prog, inst)
		}
		addrMap[s.ID] = len(vm.prog)
	}
	vm.prog = invert(vm.prog)
	vm.data = tape.String()
	return
}
