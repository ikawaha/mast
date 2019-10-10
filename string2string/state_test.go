package string2string

import (
	"fmt"
	"testing"
)

func TestEqual01(t *testing.T) {
	type pair struct {
		x *State
		y *State
	}

	s := &State{}

	testdata := []struct {
		input    pair
		expected bool
	}{
		{pair{x: s, y: s}, true},
		{pair{x: nil, y: nil}, false},
		{pair{x: nil, y: &State{}}, false},
		{pair{x: &State{}, y: nil}, false},
		{pair{&State{ID: 1}, &State{ID: 2}}, true},
		{pair{&State{IsFinal: true}, &State{IsFinal: false}}, false},
		{pair{&State{Output: map[byte]string{1: "go"}}, &State{}}, false},
		{pair{&State{Output: map[byte]string{1: "go"}}, &State{Output: map[byte]string{1: "go"}}}, true},
		{pair{&State{Output: map[byte]string{1: "go"}}, &State{Output: map[byte]string{1: "c++"}}}, false},
		{pair{&State{Output: map[byte]string{1: "go"}}, &State{Output: map[byte]string{2: "go"}}}, false},
		{pair{&State{Tail: map[string]struct{}{"go": struct{}{}}}, &State{Tail: map[string]struct{}{"go": struct{}{}}}}, true},
	}
	for _, d := range testdata {
		if got := d.input.x.Equal(d.input.y); got != d.expected {
			t.Errorf("got %v, expected %v, %v\n", got, d.expected, d)
		}
		if got := d.input.y.Equal(d.input.x); got != d.expected {
			t.Errorf("got %v, expected %v, %v\n", got, d.expected, d)
		}
	}
}

func TestEqual02(t *testing.T) {
	x := &State{ID: 1}
	y := &State{ID: 2}
	a := &State{
		Trans: map[byte]*State{1: x, 2: y},
	}
	b := &State{
		Trans: map[byte]*State{1: x, 2: y},
	}
	c := &State{
		Trans: map[byte]*State{1: y, 2: y},
	}
	d := &State{
		Trans: map[byte]*State{1: x, 2: y, 3: x},
	}
	if got, expected := a.Equal(b), true; got != expected {
		t.Errorf("got %v, expected %v\n", got, expected)
	}
	if got, expected := a.Equal(c), false; got != expected {
		t.Errorf("got %v, expected %v\n", got, expected)
	}
	if got, expected := a.Equal(c), false; got != expected {
		t.Errorf("got %v, expected %v\n", got, expected)
	}
	if got, expected := a.Equal(d), false; got != expected {
		t.Errorf("got %v, expected %v\n", got, expected)
	}

}

func TestString01(t *testing.T) {
	testdata := []struct {
		input    *State
		expected string
	}{
		{input: nil, expected: "<nil>"},
	}
	for _, d := range testdata {
		if got := d.input.String(); got != d.expected {
			t.Errorf("got %v, expected %v, %v\n", got, d.expected, d)
		}
	}
	r := &State{}
	s := State{
		ID:      1,
		Trans:   map[byte]*State{1: nil, 2: r},
		Output:  map[byte]string{3: "go", 4: "gopher"},
		Tail:    stringSet{"hello": struct{}{}},
		IsFinal: true,
		Prev:    []*State{nil, r},
	}
	fmt.Println(s.String())
}
