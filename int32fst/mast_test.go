package int32fst

import (
	"os"
	"reflect"
	"testing"
)

func TestMASTBuildMAST01(t *testing.T) {
	inp := PairSlice{}
	m := BuildMAST(inp)
	if m.StartingState.ID != 0 {
		t.Errorf("got initial State id %v, expected 0", m.StartingState.ID)
	}
	if len(m.States) != 1 {
		t.Errorf("expected: initial State only, got %v", m.States)
	}
	if len(m.FinalStates) != 0 {
		t.Errorf("expected: final State is empty, got %v", m.FinalStates)
	}
}

func TestMASTAccept01(t *testing.T) {
	input := PairSlice{
		{In: "hello", Out: 111},
		{In: "hello", Out: 222},
		{In: "111", Out: 111},
		{In: "112", Out: 112},
		{In: "112", Out: 122},
		{In: "211", Out: 345},
	}
	m := BuildMAST(input)
	for _, pair := range input {
		if ok := m.Accept(pair.In); !ok {
			t.Errorf("expected: Accept [%v]", pair.In)
		}
	}
	if ok := m.Accept("aloha"); ok {
		t.Errorf("expected: reject 'aloha'")
	}
}

func TestMASTRun01(t *testing.T) {
	input := PairSlice{
		{In: "hello", Out: 1111},
		{In: "hell", Out: 2222},
		{In: "111", Out: 111},
		{In: "112", Out: 112},
		{In: "113", Out: 122},
		{In: "211", Out: 111},
	}
	m := BuildMAST(input)
	for _, pair := range input {
		out, ok := m.Run(pair.In)
		if !ok {
			t.Errorf("expected: Accept [%v]", pair.In)
		}
		if len(out) != 1 {
			t.Errorf("input: %v, output size: got %v, expected 1", pair.In, len(out))
		}
		if out[0] != pair.Out {
			t.Errorf("input: %v, output: got %v, expected %v", pair.In, pair.Out, out[0])
		}
	}
	if out, ok := m.Run("aloha"); ok {
		t.Errorf("expected: reject 'aloha', %v", out)
	}
}

func TestMASTRun02(t *testing.T) {
	inp := PairSlice{
		{In: "hello", Out: 1111},
		{In: "hello", Out: 2222},
	}
	m := BuildMAST(inp)
	for _, pair := range inp {
		out, ok := m.Run(pair.In)
		if !ok {
			t.Errorf("expected: Accept [%v]", pair.In)
		}
		if len(out) != 2 {
			t.Errorf("input: %v, output size: got %v, expected 2", pair.In, len(out))
		}
		expected := []int32{1111, 2222}
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("input: %v, output: got %v, expected %v", pair.In, out, expected)
		}
	}
	if out, ok := m.Run("aloha"); ok {
		t.Errorf("expected: reject 'aloha', %v", out)
	}
}

func TestMASTDot01(t *testing.T) {
	input := PairSlice{
		{In: "apr", Out: 30},
		{In: "aug", Out: 31},
		{In: "dec", Out: 31},
		{In: "feb", Out: 28},
		{In: "feb", Out: 29},
	}
	m := BuildMAST(input)
	m.Dot(os.Stdout)
}

func TestMASTDot02(t *testing.T) {
	input := PairSlice{
		{In: "apr", Out: 30},
		{In: "aug", Out: 31},
		{In: "dec", Out: 31},
		{In: "feb", Out: 28},
		{In: "feb", Out: 29},
		{In: "lucene", Out: 1},
		{In: "lucid", Out: 2},
		{In: "lucifer", Out: 666},
	}
	m := BuildMAST(input)
	m.Dot(os.Stdout)
}
