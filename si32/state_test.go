package si32

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

	crs := []struct {
		call pair
		resp bool
	}{
		{pair{x: s, y: s}, true},
		{pair{x: nil, y: nil}, false},
		{pair{x: nil, y: &State{}}, false},
		{pair{x: &State{}, y: nil}, false},
		{pair{&State{ID: 1}, &State{ID: 2}}, true},
		{pair{&State{IsFinal: true}, &State{IsFinal: false}}, false},
		{pair{&State{Output: map[byte]int32{1: 555}}, &State{}}, false},
		{pair{&State{Output: map[byte]int32{1: 555}}, &State{Output: map[byte]int32{1: 555}}},
			true},
		{pair{&State{Output: map[byte]int32{1: 555}}, &State{Output: map[byte]int32{1: 444}}},
			false},
		{pair{&State{Output: map[byte]int32{1: 555}}, &State{Output: map[byte]int32{2: 555}}},
			false},
		{pair{&State{Tail: map[int32]bool{555: true}}, &State{Tail: map[int32]bool{555: true}}}, true},
	}
	for _, cr := range crs {
		if rst := cr.call.x.Equal(cr.call.y); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
		if rst := cr.call.y.Equal(cr.call.x); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
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
	if rst, exp := a.Equal(b), true; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}

	a.SetOutput('a', 1)
	b.SetOutput('a', 2)
	if rst, exp := a.Equal(b), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}

	if rst, exp := a.Equal(c), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.Equal(c), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.Equal(d), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}

}

func TestStateString01(t *testing.T) {
	crs := []struct {
		call *State
		resp string
	}{
		{nil, "<nil>"},
	}
	for _, cr := range crs {
		if rst := cr.call.String(); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
	}
	r := &State{}
	s := State{
		ID:      1,
		Trans:   map[byte]*State{1: nil, 2: r},
		Output:  map[byte]int32{3: 555, 4: 888},
		Tail:    int32Set{1111: true},
		IsFinal: true,
	}
	fmt.Println(s.String())
}
