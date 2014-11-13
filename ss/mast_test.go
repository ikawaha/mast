package ss

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestBuildMast01(t *testing.T) {
	inp := PairSlice{}
	m := buildMast(inp)
	if m.initialState.ID != 0 {
		t.Errorf("got initial state id %v, expected 0\n", m.initialState.ID)
	}
	if len(m.states) != 1 {
		t.Errorf("expected: initial state only, got %v\n", m.states)
	}
	if len(m.finalStates) != 0 {
		t.Errorf("expected: final state is empty, got %v\n", m.finalStates)
	}
}

func TestAccept01(t *testing.T) {
	inp := PairSlice{
		{"hello", "world"},
		{"hello", "goodby"},
		{"111", "aaa"},
		{"112", "aab"},
		{"112", "abb"},
		{"211", "cde"},
	}
	m := buildMast(inp)
	for _, pair := range inp {
		if ok := m.accept(pair.In); !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
	}
	if ok := m.accept("aloha"); ok {
		t.Errorf("expected: reject \"aloha\"\n")
	}
}

func TestRun01(t *testing.T) {
	inp := PairSlice{
		{"hello", "world"},
		{"hell", "daemon"},
		{"111", "aaa"},
		{"112", "aab"},
		{"113", "abb"},
		{"211", "aaa"},
	}
	m := buildMast(inp)
	for _, pair := range inp {
		out, ok := m.run(pair.In)
		if !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
		if len(out) != 1 {
			t.Errorf("input: %v, output size: got %v, expected 1\n", pair.In, len(out))
		}
		if out[0] != pair.Out {
			t.Errorf("input: %v, output: got %v, expected %v\n", pair.In, pair.Out, out[0])
		}
	}
	if out, ok := m.run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v\n", out)
	}
}

func TestMastRun02(t *testing.T) {
	inp := PairSlice{
		{"hello", "world"},
		{"hello", "goodby"},
	}
	m := buildMast(inp)
	for _, pair := range inp {
		out, ok := m.run(pair.In)
		if !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
		if len(out) != 2 {
			t.Errorf("input: %v, output size: got %v, expected 2\n", pair.In, len(out))
		}
		expected := []string{"world", "goodby"}
		sort.Strings(out)
		sort.Strings(expected)
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("input: %v, output: got %v, expected %v\n", pair.In, out, expected)
		}
	}
	if out, ok := m.run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v\n", out)
	}
}

func TestMastDot01(t *testing.T) {
	inp := PairSlice{
		{"apr", "30"},
		{"aug", "31"},
		{"dec", "31"},
		{"feb", "28"},
		{"feb", "29"},
	}
	m := buildMast(inp)
	m.dot(os.Stdout)
}

func TestToBytes01(t *testing.T) {
	cr := []struct {
		in  int
		out []byte
	}{
		{7, []byte{7}},
		{255, []byte{255}},
		{256, []byte{0, 1}},
		{512, []byte{0, 2}},
		{1024, []byte{0, 4}},
		{0xFFFFFF, []byte{255, 255, 255}},
	}
	for _, s := range cr {
		if r := toBytes(s.in); !reflect.DeepEqual(s.out, r) {
			t.Errorf("got %v, expected %v\n", r, s.out)
		}
	}
}

func TestCompile01(t *testing.T) {
	inp := PairSlice{
		{"1a22xss", "world"},
		{"1b22yss", "goodby"},
	}
	m := buildMast(inp)
	_, e := m.compile()
	if e != nil {
		t.Errorf("unexpected error: %v", e)
	}
}

func TestMastCompile02(t *testing.T) {
	inp := PairSlice{
		{"abc", "123"},
		{"abc", "456"},
	}
	m := buildMast(inp)
	vm, e := m.compile()
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	outs := vm.Search("abc")
	exp := []string{"123", "456"}
	sort.Strings(outs)
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("got %v, expected %v\n", outs, exp)
	}
}
