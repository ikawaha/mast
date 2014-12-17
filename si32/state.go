package si32

import "fmt"

type int32Set map[int32]bool

type state struct {
	ID      int
	Trans   map[byte]*state
	Output  map[byte]int32
	Tail    int32Set
	IsFinal bool
	Prev    []*state
	hcode   int64
}

func newState() (n *state) {
	n = new(state)
	n.Trans = make(map[byte]*state)
	n.Output = make(map[byte]int32)
	n.Tail = make(int32Set)
	return
}

func (n *state) hasTail() bool {
	return len(n.Tail) != 0
}

func (n *state) addTail(t int32) {
	n.Tail[t] = true
}

/*
func (n *state) setTail(s int32Set) {
	n.Tail = s
}
*/

func (n *state) tails() (t []int32) {
	t = make([]int32, 0, len(n.Tail))
	for item := range n.Tail {
		t = append(t, item)
	}
	return
}

func (n *state) removeOutput(ch byte) {
	const magic = 8191
	if out, ok := n.Output[ch]; ok && out != 0 {
		n.hcode -= (int64(ch) + int64(out)) * magic
	}
	delete(n.Output, ch)
}

func (n *state) setOutput(ch byte, out int32) {
	if out == 0 {
		return
	}
	n.Output[ch] = out

	const magic = 8191
	n.hcode += (int64(ch) + int64(out)) * magic
}

func (n *state) setTransition(ch byte, next *state) {
	n.Trans[ch] = next

	const magic = 1001
	n.hcode += (int64(ch) + int64(next.ID)) * magic
}

func (n *state) setInvTransition() {
	for _, next := range n.Trans {
		next.Prev = append(next.Prev, n)
	}
}

func (n *state) renew() {
	n.Trans = make(map[byte]*state)
	n.Output = make(map[byte]int32)
	n.Tail = make(int32Set)
	n.IsFinal = false
	n.Prev = make([]*state, 0)
	n.hcode = 0
}

func (n *state) eq(dst *state) bool {
	if n == nil || dst == nil {
		return false
	}
	if n == dst {
		return true
	}
	if n.hcode != dst.hcode {
		return false
	}
	if len(n.Trans) != len(dst.Trans) ||
		len(n.Output) != len(dst.Output) ||
		len(n.Tail) != len(dst.Tail) ||
		n.IsFinal != dst.IsFinal {
		return false
	}
	for ch, next := range n.Trans {
		if dst.Trans[ch] != next {
			return false
		}
	}
	for ch, out := range n.Output {
		if dst.Output[ch] != out {
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
	for ch := range n.Trans {
		ret += fmt.Sprintf("%X02/%v -->%p, ", ch, n.Output[ch], n.Trans[ch])
	}
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
