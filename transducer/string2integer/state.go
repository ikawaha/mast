package string2integer

import (
	"fmt"
)

const (
	magic = 8191
)

type intSet map[int]struct{}

type state struct {
	ID         int
	transition map[byte]*state
	outputs    map[byte]int
	tail       intSet
	isFinal    bool
	prev       []*state
	hashCode   uint
}

func newState() *state {
	return &state{
		transition: map[byte]*state{},
		outputs:    map[byte]int{},
		tail:       intSet{},
	}
}

func (s *state) hasTail() bool {
	return len(s.tail) != 0
}

func (s *state) addTail(v int) {
	s.tail[v] = struct{}{}
}

func (s *state) tails() []int {
	ret := make([]int, 0, len(s.tail))
	for item := range s.tail {
		ret = append(ret, item)
	}
	return ret
}

func (s *state) deleteOutput(ch byte) {
	delete(s.outputs, ch)
}

func (s *state) setOutput(ch byte, out int) {
	s.outputs[ch] = out
}

func (s *state) setTransition(ch byte, next *state) {
	s.transition[ch] = next

	s.hashCode += (uint(ch) + uint(next.ID)) * magic
}

func (s *state) setInvTransition() {
	for _, next := range s.transition {
		next.prev = append(next.prev, s)
	}
}

func (s *state) clear() {
	s.transition = make(map[byte]*state)
	s.outputs = make(map[byte]int)
	s.tail = intSet{}
	s.isFinal = false
	s.prev = []*state{}
	s.hashCode = 0
}

func (s *state) equal(dst *state) bool {
	if s == nil || dst == nil {
		return false
	}
	if s == dst {
		return true
	}
	if s.hashCode != dst.hashCode {
		return false
	}
	if len(s.transition) != len(dst.transition) ||
		len(s.outputs) != len(dst.outputs) ||
		len(s.tail) != len(dst.tail) ||
		s.isFinal != dst.isFinal {
		return false
	}
	for ch, next := range s.transition {
		if dst.transition[ch] != next {
			return false
		}
	}
	for ch, out := range s.outputs {
		if dst.outputs[ch] != out {
			return false
		}
	}
	for item := range s.tail {
		if _, ok := dst.tail[item]; !ok {
			return false
		}
	}
	return true
}

// String returns a string representaion of a node for debug.
func (s *state) String() string {
	ret := ""
	if s == nil {
		return "<nil>"
	}
	ret += fmt.Sprintf("%d[%p]:", s.ID, s)
	for ch := range s.transition {
		ret += fmt.Sprintf("%X/%v -->%p, ", ch, s.outputs[ch], s.transition[ch])
	}
	if s.isFinal {
		ret += fmt.Sprintf(" (tail:%v) ", s.tails())
	}
	ret += fmt.Sprint("<--(")
	for _, s := range s.prev {
		ret += fmt.Sprintf("%p, ", s)
	}
	ret += fmt.Sprint(")")
	return ret
}
