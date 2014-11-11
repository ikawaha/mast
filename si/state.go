package si

import (
	"fmt"
)

type intSet map[int]bool

type state struct {
	ID      int
	Trans   map[byte]*state
	Tail    intSet
	IsFinal bool
	Prev    []*state
}

func newState() (n *state) {
	n = new(state)
	n.Trans = make(map[byte]*state)
	n.Tail = make(intSet)
	return
}

func (n *state) hasTail() bool {
	return len(n.Tail) != 0
}

func (n *state) addTail(t int) {
	n.Tail[t] = true
}

func (n *state) setTail(s intSet) {
	n.Tail = s
}

func (n *state) tails() (t []int) {
	t = make([]int, 0, len(n.Tail))
	for item := range n.Tail {
		t = append(t, item)
	}
	return
}

func (n *state) setTransition(ch byte, next *state) {
	n.Trans[ch] = next
}

func (n *state) setInvTransition() {
	for _, next := range n.Trans {
		next.Prev = append(next.Prev, n)
	}
}

func (n *state) renew() {
	n.Trans = make(map[byte]*state)
	n.Tail = make(intSet)
	n.IsFinal = false
	n.Prev = make([]*state, 0)
}

func (n *state) eq(dst *state) bool {
	if n == nil || dst == nil {
		return false
	}
	if n == dst {
		return true
	}
	if len(n.Trans) != len(dst.Trans) ||
		len(n.Tail) != len(dst.Tail) ||
		n.IsFinal != dst.IsFinal {
		return false
	}
	for ch, next := range n.Trans {
		if dst.Trans[ch] != next {
			return false
		}
	}
	for item := range n.Tail {
		if !dst.Tail[item] {
			return false
		}
	}
	return true
}

// String returns a string representaion of a node for debug.
func (n *state) String() string {
	ret := ""
	if n == nil {
		return "<nil>"
	}
	ret += fmt.Sprintf("%d[%p]:", n.ID, n)
	if n.IsFinal {
		ret += fmt.Sprintf(" (tail:%v) ", n.tails())
	}
	ret += fmt.Sprint("<--(")
	for _, s := range n.Prev {
		ret += fmt.Sprintf("%p, ", s)
	}
	ret += fmt.Sprint(")")
	return ret
}
