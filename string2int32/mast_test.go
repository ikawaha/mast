package string2int32

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestMASTBuildMAST01(t *testing.T) {
	inp := PairSlice{}
	m := BuildMAST(inp)
	if m.StartingState.ID != 0 {
		t.Errorf("got initial State id %v, expected 0\n", m.StartingState.ID)
	}
	if len(m.States) != 1 {
		t.Errorf("expected: initial State only, got %v\n", m.States)
	}
	if len(m.FinalStates) != 0 {
		t.Errorf("expected: final State is empty, got %v\n", m.FinalStates)
	}
}

func TestMASTAccept01(t *testing.T) {
	inp := PairSlice{
		{"hello", 111},
		{"hello", 222},
		{"111", 111},
		{"112", 112},
		{"112", 122},
		{"211", 345},
	}
	m := BuildMAST(inp)
	for _, pair := range inp {
		if ok := m.Accept(pair.In); !ok {
			t.Errorf("expected: Accept [%v]\n", pair.In)
		}
	}
	if ok := m.Accept("aloha"); ok {
		t.Errorf("expected: reject \"aloha\"\n")
	}
}

func TestMASTRun01(t *testing.T) {
	inp := PairSlice{
		{"hello", 1111},
		{"hell", 2222},
		{"111", 111},
		{"112", 112},
		{"113", 122},
		{"211", 111},
	}
	m := BuildMAST(inp)
	for _, pair := range inp {
		out, ok := m.Run(pair.In)
		if !ok {
			t.Errorf("expected: Accept [%v]\n", pair.In)
		}
		if len(out) != 1 {
			t.Errorf("input: %v, output size: got %v, expected 1\n", pair.In, len(out))
		}
		if out[0] != pair.Out {
			t.Errorf("input: %v, output: got %v, expected %v\n", pair.In, pair.Out, out[0])
		}
	}
	if out, ok := m.Run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v\n", out)
	}
}

func TestMASTRun02(t *testing.T) {
	inp := PairSlice{
		{"hello", 1111},
		{"hello", 2222},
	}
	m := BuildMAST(inp)
	for _, pair := range inp {
		out, ok := m.Run(pair.In)
		if !ok {
			t.Errorf("expected: Accept [%v]\n", pair.In)
		}
		if len(out) != 2 {
			t.Errorf("input: %v, output size: got %v, expected 2\n", pair.In, len(out))
		}
		expected := []int32{1111, 2222}
		sort.Sort(int32Slice(out))
		sort.Sort(int32Slice(expected))
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("input: %v, output: got %v, expected %v\n", pair.In, out, expected)
		}
	}
	if out, ok := m.Run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v\n", out)
	}
}

func TestMASTDot01(t *testing.T) {
	inp := PairSlice{
		{"apr", 30},
		{"aug", 31},
		{"dec", 31},
		{"feb", 28},
		{"feb", 29},
	}
	m := BuildMAST(inp)
	m.Dot(os.Stdout)
}

func TestMASTDot02(t *testing.T) {
	inp := PairSlice{
		{"apr", 30},
		{"aug", 31},
		{"dec", 31},
		{"feb", 28},
		{"feb", 29},
		{"lucene", 1},
		{"lucid", 2},
		{"lucifer", 666},
	}
	m := BuildMAST(inp)
	m.Dot(os.Stdout)
}
