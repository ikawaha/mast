package stringfst

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestBuildMast01(t *testing.T) {
	input := PairSlice{}
	m := BuildMAST(input)
	if m.StartingState.ID != 0 {
		t.Errorf("got initial state id %v, expected 0", m.StartingState.ID)
	}
	if len(m.States) != 1 {
		t.Errorf("expected initial state only, got %v", m.States)
	}
	if len(m.FinalStates) != 0 {
		t.Errorf("expected final state is empty, got %v", m.FinalStates)
	}
}

func TestAccept01(t *testing.T) {
	input := PairSlice{
		{In: "hello", Out: "world"},
		{In: "hello", Out: "goodbye"},
		{In: "111", Out: "aaa"},
		{In: "112", Out: "aab"},
		{In: "112", Out: "abb"},
		{In: "211", Out: "cde"},
	}
	m := BuildMAST(input)
	for _, pair := range input {
		if ok := m.Accept(pair.In); !ok {
			t.Errorf("expected: accept '%v'", pair.In)
		}
	}
	if ok := m.Accept("aloha"); ok {
		t.Errorf("expected: reject 'aloha'")
	}
}

func TestRun01(t *testing.T) {
	input := PairSlice{
		{In: "hello", Out: "world"},
		{In: "hell", Out: "daemon"},
		{In: "111", Out: "aaa"},
		{In: "112", Out: "aab"},
		{In: "113", Out: "abb"},
		{In: "211", Out: "aaa"},
	}
	m := BuildMAST(input)
	for _, pair := range input {
		out, ok := m.Run(pair.In)
		if !ok {
			t.Errorf("expected: accept [%v]", pair.In)
		}
		if len(out) != 1 {
			t.Errorf("input %v, output size: got %v, expected 1", pair.In, len(out))
		}
		if out[0] != pair.Out {
			t.Errorf("input %v, output: got %v, expected %v", pair.In, pair.Out, out[0])
		}
	}
	if out, ok := m.Run("aloha"); ok {
		t.Errorf("expected: reject 'aloha', %v", out)
	}
}

func TestMastRun02(t *testing.T) {
	input := PairSlice{
		{In: "hello", Out: "world"},
		{In: "hello", Out: "goodbye"},
	}
	m := BuildMAST(input)
	for _, pair := range input {
		out, ok := m.Run(pair.In)
		if !ok {
			t.Errorf("expected: accept [%v]", pair.In)
		}
		if len(out) != 2 {
			t.Errorf("input: %v, output size: got %v, expected 2", pair.In, len(out))
		}
		expected := []string{"world", "goodbye"}
		sort.Strings(out)
		sort.Strings(expected)
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("input: %v, output: got %v, expected %v", pair.In, out, expected)
		}
	}
	if out, ok := m.Run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v", out)
	}
}

func TestMastDot01(t *testing.T) {
	input := PairSlice{
		{In: "apr", Out: "30"},
		{In: "aug", Out: "31"},
		{In: "dec", Out: "31"},
		{In: "feb", Out: "28"},
		{In: "feb", Out: "29"},
	}
	m := BuildMAST(input)
	m.Dot(os.Stdout)
}
