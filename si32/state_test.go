package si32

import "testing"
import "fmt"

func TestEq01(t *testing.T) {
	type pair struct {
		x *state
		y *state
	}

	s := &state{}

	crs := []struct {
		call pair
		resp bool
	}{
		{pair{x: s, y: s}, true},
		{pair{x: nil, y: nil}, false},
		{pair{x: nil, y: &state{}}, false},
		{pair{x: &state{}, y: nil}, false},
		{pair{&state{ID: 1}, &state{ID: 2}}, true},
		{pair{&state{IsFinal: true}, &state{IsFinal: false}}, false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{}}, false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{1: 555}}},
			true},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{1: 444}}},
			false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{2: 555}}},
			false},
		{pair{&state{Tail: map[int32]bool{555: true}}, &state{Tail: map[int32]bool{555: true}}}, true},
	}
	for _, cr := range crs {
		if rst := cr.call.x.eq(cr.call.y); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
		if rst := cr.call.y.eq(cr.call.x); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
	}
}

func TestEq02(t *testing.T) {
	x := &state{ID: 1}
	y := &state{ID: 2}
	a := &state{
		Trans: map[byte]*state{1: x, 2: y},
	}
	b := &state{
		Trans: map[byte]*state{1: x, 2: y},
	}
	c := &state{
		Trans: map[byte]*state{1: y, 2: y},
	}
	d := &state{
		Trans: map[byte]*state{1: x, 2: y, 3: x},
	}
	if rst, exp := a.eq(b), true; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.eq(c), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.eq(c), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.eq(d), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}

}

func TestString01(t *testing.T) {
	crs := []struct {
		call *state
		resp string
	}{
		{nil, "<nil>"},
	}
	for _, cr := range crs {
		if rst := cr.call.String(); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
	}
	r := &state{}
	s := state{
		ID:      1,
		Trans:   map[byte]*state{1: nil, 2: r},
		Output:  map[byte]int32{3: 555, 4: 888},
		Tail:    int32Set{1111: true},
		IsFinal: true,
//		Prev:    []*state{nil, r},
	}
	fmt.Println(s.String())
}
