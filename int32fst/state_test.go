package int32fst

import (
	"fmt"
	"testing"
)

func TestStateEq01(t *testing.T) {
	type pair struct {
		x *State
		y *State
	}

	s := &State{}

	var (
		testdata = []struct {
			input    pair
			expected bool
		}{
			{input: pair{x: s, y: s}, expected: true},
			{input: pair{x: nil, y: nil}, expected: false},
			{input: pair{x: nil, y: &State{}}},
			{input: pair{x: &State{}, y: nil}},
			{input: pair{&State{ID: 1}, &State{ID: 2}}, expected: true},
			{input: pair{&State{IsFinal: true}, &State{IsFinal: false}}},
			{input: pair{&State{Output: map[byte]int32{1: 555}}, &State{}}},
			{input: pair{&State{Output: map[byte]int32{1: 555}}, &State{Output: map[byte]int32{1: 555}}}, expected: true},
			{input: pair{&State{Output: map[byte]int32{1: 555}}, &State{Output: map[byte]int32{1: 444}}}},
			{input: pair{&State{Output: map[byte]int32{1: 555}}, &State{Output: map[byte]int32{2: 555}}}},
			{input: pair{&State{Tail: map[int32]struct{}{555: struct{}{}}}, &State{Tail: map[int32]struct{}{555: struct{}{}}}}, expected: true},
		}
	)
	for _, d := range testdata {
		if got := d.input.x.Equal(d.input.y); got != d.expected {
			t.Errorf("got %v, expected %v, %v", got, d.expected, d)
		}
		if rst := d.input.y.Equal(d.input.x); rst != d.expected {
			t.Errorf("got %v, expected %v, %v", rst, d.expected, d)
		}
	}
}

func TestStateEq02(t *testing.T) {
	x := &State{ID: 1}
	y := &State{ID: 2}
	a := &State{
		Trans:  map[byte]*State{1: x, 2: y},
		Output: make(map[byte]int32),
	}
	b := &State{
		Trans:  map[byte]*State{1: x, 2: y},
		Output: make(map[byte]int32),
	}
	c := &State{
		Trans: map[byte]*State{1: y, 2: y},
	}
	d := &State{
		Trans: map[byte]*State{1: x, 2: y, 3: x},
	}
	if got, expected := a.Equal(b), true; got != expected {
		t.Errorf("got %v, expected %v", got, expected)
	}

	a.SetOutput('a', 1)
	b.SetOutput('a', 2)
	if got, expected := a.Equal(b), false; got != expected {
		t.Errorf("got %v, expected %v", got, expected)
	}

	if got, expected := a.Equal(c), false; got != expected {
		t.Errorf("got %v, expected %v", got, expected)
	}
	if got, expected := a.Equal(c), false; got != expected {
		t.Errorf("got %v, expected %v", got, expected)
	}
	if got, expected := a.Equal(d), false; got != expected {
		t.Errorf("got %v, expected %v", got, expected)
	}

}

func TestStateString(t *testing.T) {
	expected := "<nil>"
	var s *State
	if got := s.String(); got != "<nil>" {
		t.Errorf("got %v, expected %v", got, expected)
	}
	r := &State{}
	s = &State{
		ID:      1,
		Trans:   map[byte]*State{1: nil, 2: r},
		Output:  map[byte]int32{3: 123, 4: 456},
		Tail:    int32Set{789: struct{}{}},
		IsFinal: true,
	}
	fmt.Println(s.String())

}
