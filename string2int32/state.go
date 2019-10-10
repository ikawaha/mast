package string2int32

import "fmt"

type int32Set map[int32]struct{}

// State represents a state of automata.
type State struct {
	ID      int
	Trans   map[byte]*State
	Output  map[byte]int32
	Tail    int32Set
	IsFinal bool
	hcode   int64
}

// NewState constructs a new state.
func NewState() *State {
	return &State{
		Trans:  map[byte]*State{},
		Output: map[byte]int32{},
		Tail:   int32Set{},
	}
}

// HasTail returns true if the state has tail items.
func (s *State) HasTail() bool {
	return len(s.Tail) != 0
}

// AddTail adds an item to the tail set.
func (s *State) AddTail(t int32) {
	s.Tail[t] = struct{}{}
}

// Tails returns an array of items of the tail.
func (s *State) Tails() []int32 {
	ret := make([]int32, 0, len(s.Tail))
	for item := range s.Tail {
		ret = append(ret, item)
	}
	return ret
}

// RemoveOutput removes the output associated with the transition at the given character.
func (s *State) RemoveOutput(ch byte) {
	const magic = 8191
	if out, ok := s.Output[ch]; ok && out != 0 {
		s.hcode -= (int64(ch) + int64(out)) * magic
	}
	delete(s.Output, ch)
}

// SetOutput sets the output associated with the transition at the given character.
func (s *State) SetOutput(ch byte, out int32) {
	s.Output[ch] = out

	const magic = 8191
	s.hcode += (int64(ch) + int64(out)) * magic
}

// SetTransition sets the transition associated with the given character.
func (s *State) SetTransition(ch byte, next *State) {
	s.Trans[ch] = next

	const magic = 1001
	s.hcode += (int64(ch) + int64(next.ID)) * magic
}

// Clear clears the state.
func (s *State) Clear() {
	s.Trans = make(map[byte]*State)
	s.Output = make(map[byte]int32)
	s.Tail = make(int32Set)
	s.IsFinal = false
	s.hcode = 0
}

// Equal returns whether two states are equal.
func (s *State) Equal(dst *State) bool {
	if s == nil || dst == nil {
		return false
	}
	if s == dst {
		return true
	}
	if s.hcode != dst.hcode {
		return false
	}
	if len(s.Trans) != len(dst.Trans) ||
		len(s.Output) != len(dst.Output) ||
		len(s.Tail) != len(dst.Tail) ||
		s.IsFinal != dst.IsFinal {
		return false
	}
	for ch, next := range s.Trans {
		if dst.Trans[ch] != next {
			return false
		}
	}
	for ch, out := range s.Output {
		if dst.Output[ch] != out {
			return false
		}
	}
	for item := range s.Tail {
		if _, ok := dst.Tail[item]; !ok {
			return false
		}
	}
	return true
}

// String returns the string representation of the state.
func (s *State) String() string {
	ret := ""
	if s == nil {
		return "<nil>"
	}
	ret += fmt.Sprintf("%d[%p]:", s.ID, s)
	for ch := range s.Trans {
		ret += fmt.Sprintf("%X02/%v -->%p, ", ch, s.Output[ch], s.Trans[ch])
	}
	if s.IsFinal {
		ret += fmt.Sprintf(" (tail:%v) ", s.Tails())
	}
	return ret
}
