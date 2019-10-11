package string2int32

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestFSTRun01(t *testing.T) {
	input := PairSlice{
		{In: "feb", Out: 28},
		{In: "feb", Out: 29},
		{In: "feb", Out: 30},
		{In: "dec", Out: 31},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	t.Run("graph output", func(t *testing.T) {
		var b bytes.Buffer
		m.Dot(&b)
		expected := `digraph G {
	rankdir=LR;
	node [shape=circle]
	0 [peripheries = 2];
	3 [peripheries = 2];
	1 -> 0 [label="63(c)/ε"];
	2 -> 1 [label="65(e)/ε"];
	4 -> 3 [label="62(b)/ε [28 29 30]"];
	5 -> 4 [label="65(e)/ε"];
	6 -> 2 [label="64(d)/31"];
	6 -> 5 [label="66(f)/ε"];
}
`
		if got := b.String(); got != expected {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("machine code", func(t *testing.T) {
		expected := `  0 OUTPUT	64(d) 7
  1 [31]
  2 MATCHB	66(f) 1
  3 MATCHB	65(e) 1
  4 MATCHB	62(b) 1
  5 ACCEPTB	1 0
  6 [3]
  7 [0] [28 29 30]
  8 MATCHB	65(e) 1
  9 MATCHB	63(c) 1
 10 ACCEPTB	0 0
`
		if got := fmt.Sprint(fst); got != expected {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("running", func(t *testing.T) {
		expected := []Configuration{
			{
				PC:      5,
				Head:    3,
				Outputs: []int32{28, 29, 30},
			},
		}
		got, ok := fst.Run("feb")
		if !ok {
			t.Errorf("input:feb, config:%v, Accept:%v", got, ok)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})
}

func TestFSTRun02(t *testing.T) {
	inp := PairSlice{
		{In: "feb", Out: 28},
		{In: "feb", Out: 29},
		{In: "feb", Out: 30},
		{In: "dec", Out: 31},
	}
	m := BuildMAST(inp)
	fst, err := m.BuildFST()
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	t.Run("graph output", func(t *testing.T) {
		var b bytes.Buffer
		m.Dot(&b)
		expected := `digraph G {
	rankdir=LR;
	node [shape=circle]
	0 [peripheries = 2];
	3 [peripheries = 2];
	1 -> 0 [label="63(c)/ε"];
	2 -> 1 [label="65(e)/ε"];
	4 -> 3 [label="62(b)/ε [28 29 30]"];
	5 -> 4 [label="65(e)/ε"];
	6 -> 2 [label="64(d)/31"];
	6 -> 5 [label="66(f)/ε"];
}
`
		if got := b.String(); got != expected {
			fmt.Println(fst.Data)
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("machine code", func(t *testing.T) {
		expected := `  0 OUTPUT	64(d) 7
  1 [31]
  2 MATCHB	66(f) 1
  3 MATCHB	65(e) 1
  4 MATCHB	62(b) 1
  5 ACCEPTB	1 0
  6 [3]
  7 [0] [28 29 30]
  8 MATCHB	65(e) 1
  9 MATCHB	63(c) 1
 10 ACCEPTB	0 0
`
		if got := fmt.Sprint(fst); got != expected {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("running", func(t *testing.T) {
		expected := []Configuration{
			{
				PC:      10,
				Head:    3,
				Outputs: []int32{31},
			},
		}
		got, ok := fst.Run("dec")
		if !ok {
			t.Errorf("input:dec, config:%v, Accept:%v", got, ok)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})
}

func TestFSTRun03(t *testing.T) {
	inp := PairSlice{
		{In: "feb", Out: 0},
		{In: "february", Out: 1},
	}
	m := BuildMAST(inp)
	fst, err := m.BuildFST()
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	t.Run("graph output", func(t *testing.T) {
		var b bytes.Buffer
		m.Dot(&b)
		expected := `digraph G {
	rankdir=LR;
	node [shape=circle]
	0 [peripheries = 2];
	5 [peripheries = 2];
	1 -> 0 [label="79(y)/ε"];
	2 -> 1 [label="72(r)/ε"];
	3 -> 2 [label="61(a)/ε"];
	4 -> 3 [label="75(u)/ε"];
	5 -> 4 [label="72(r)/1"];
	6 -> 5 [label="62(b)/ε"];
	7 -> 6 [label="65(e)/ε"];
	8 -> 7 [label="66(f)/ε"];
}
`
		if got := b.String(); got != expected {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("machine code", func(t *testing.T) {
		expected := `  0 MATCHB	66(f) 1
  1 MATCHB	65(e) 1
  2 MATCHB	62(b) 1
  3 ACCEPT	0 0
  4 OUTPUTB	72(r) 1
  5 [1]
  6 MATCHB	75(u) 1
  7 MATCHB	61(a) 1
  8 MATCHB	72(r) 1
  9 MATCHB	79(y) 1
 10 ACCEPTB	0 0
`
		if got := fmt.Sprint(fst); got != expected {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("running", func(t *testing.T) {
		expected := []Configuration{
			{
				PC:      3,
				Head:    3,
				Outputs: []int32{0},
			},
			{
				PC:      10,
				Head:    8,
				Outputs: []int32{1},
			},
		}
		got, ok := fst.Run("february")
		if !ok {
			t.Errorf("input:dec, config:%v, Accept:%v", got, ok)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})
}

func TestFSTSearch01(t *testing.T) {
	input := PairSlice{
		{In: "1a22xss", Out: 111},
		{In: "1a22", Out: 111},
		{In: "1b22yss", Out: 222},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	for _, p := range input {
		outs := fst.Search(p.In)
		if !reflect.DeepEqual(outs, []int32{p.Out}) {
			t.Errorf("input %v, got %v, expected %v\n", p.In, outs, []int32{p.Out})
		}
	}
}

func TestFSTSearch02(t *testing.T) {
	input := PairSlice{
		{In: "hell", Out: 666},
		{In: "hello", Out: 111},
		{In: "goodbye", Out: 222},
		{In: "goodbye", Out: 333},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	testdata := []struct {
		input    string
		expected []int32
	}{
		{input: "hell", expected: []int32{666}},
		{input: "hello", expected: []int32{111}},
		{input: "goodbye", expected: []int32{222, 333}},
	}
	for _, d := range testdata {
		outs := fst.Search(d.input)
		if !reflect.DeepEqual(outs, d.expected) {
			t.Errorf("input %v, got %v, expected %v\n", d.input, outs, d.expected)
		}
	}
}

func TestFSTSearch03(t *testing.T) {
	input := PairSlice{
		{In: "hell", Out: 0},
		{In: "hello", Out: 0},
		{In: "goodbye", Out: 0},
		{In: "goodbye", Out: 0},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	testdata := []struct {
		input    string
		expected []int32
	}{
		{input: "hell", expected: []int32{0}},
		{input: "hello", expected: []int32{0}},
		{input: "goodbye", expected: []int32{0}},
	}
	for _, d := range testdata {
		outs := fst.Search(d.input)
		if !reflect.DeepEqual(outs, d.expected) {
			t.Errorf("input %v, got %v, expected %v\n", d.input, outs, d.expected)
		}
	}
}

func TestFSTSearch06(t *testing.T) {
	input := PairSlice{
		{In: "こんにちは", Out: 111},
		{In: "世界", Out: 222},
		{In: "すもももももも", Out: 333},
		{In: "すもも", Out: 333},
		{In: "すもも", Out: 444},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	testdata := []struct {
		input    string
		expected []int32
	}{
		{input: "すもも", expected: []int32{333, 444}},
		{input: "こんにちわ"},
		{input: "こんにちは", expected: []int32{111}},
		{input: "世界", expected: []int32{222}},
		{input: "すもももももも", expected: []int32{333}},
		{input: "すももももももも", expected: nil},
		{input: "すも", expected: nil},
		{input: "すもう", expected: nil},
	}

	for _, d := range testdata {
		outs := fst.Search(d.input)
		if !reflect.DeepEqual(outs, d.expected) {
			t.Errorf("input:%v, got %v, expected %v\n", d.input, outs, d.expected)
		}
	}
}

func TestFSTPrefixSearch01(t *testing.T) {
	input := PairSlice{
		{In: "こんにちは", Out: 111},
		{In: "世界", Out: 222},
		{In: "すもももももも", Out: 333},
		{In: "すもも", Out: 333},
		{In: "すもも", Out: 444},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	testdata := []struct {
		input   string
		pos     int
		outputs []int32
	}{
		{input: "すもも", pos: 9, outputs: []int32{333, 444}},
		{input: "こんにちわ", pos: -1},
		{input: "こんにちは", pos: 15, outputs: []int32{111}},
		{input: "世界", pos: 6, outputs: []int32{222}},
		{input: "すもももももも", pos: 21, outputs: []int32{333}},
		{input: "すもももももももものうち", pos: 21, outputs: []int32{333}},
		{input: "すも", pos: -1},
		{input: "すもう", pos: -1},
	}

	for _, d := range testdata {
		pos, outs := fst.PrefixSearch(d.input)
		if !reflect.DeepEqual(outs, d.outputs) {
			t.Errorf("input:%v, got %v, expected %v\n", d.input, outs, d.outputs)
		}
		if pos != d.pos {
			t.Errorf("input:%v, got %v, expected %v\n", d.input, pos, d.pos)
		}
	}
}

func TestFSTCommonPrefixSearch01(t *testing.T) {
	input := PairSlice{
		{"こんにちは", 111},
		{"世界", 222},
		{"すもももももも", 333},
		{"すもも", 333},
		{"すもも", 444},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	testdata := []struct {
		input   string
		lens    []int
		outputs [][]int32
	}{
		{input: "すもも", lens: []int{9}, outputs: [][]int32{{333, 444}}},
		{input: "こんにちわ"},
		{input: "こんにちは", lens: []int{15}, outputs: [][]int32{{111}}},
		{input: "世界", lens: []int{6}, outputs: [][]int32{{222}}},
		{input: "すもももももも", lens: []int{9, 21}, outputs: [][]int32{{333, 444}, {333}}},
		{input: "すもももももももものうち", lens: []int{9, 21}, outputs: [][]int32{{333, 444}, {333}}},
		{input: "すも", lens: nil, outputs: nil},
		{input: "すもう", lens: nil, outputs: nil},
	}

	for _, d := range testdata {
		lens, outs := fst.CommonPrefixSearch(d.input)
		if !reflect.DeepEqual(lens, d.lens) {
			t.Errorf("input:%v, got %v %v, expected %v %v\n", d.input, lens, outs, d.lens, d.outputs)
		}
		for i := range lens {
			sort.Sort(int32Slice(outs[i]))
			sort.Sort(int32Slice(d.outputs[i]))
			if !reflect.DeepEqual(outs[i], d.outputs[i]) {
				t.Errorf("input:%v, got %v %v, expected %v %v\n", d.input, lens, outs, d.lens, d.outputs)
			}
		}
	}
}

func TestFSTSaveAndLoad01(t *testing.T) {
	inp := PairSlice{
		{"feb", 28},
		{"feb", 29},
		{"apr", 30},
		{"jan", 31},
		{"jun", 30},
		{"jul", 31},
		{"dec", 31},
	}

	m := BuildMAST(inp)
	org, e := m.BuildFST()
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}

	var b bytes.Buffer
	n, e := org.WriteTo(&b)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	if n != int64(b.Len()) {
		t.Errorf("write len: got %v, expected %v", n, b.Len())
	}

	rst, e := Read(&b)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}

	if !reflect.DeepEqual(org.Data, rst.Data) {
		t.Errorf("Data:got %v, expected %v\n", rst.Data, org.Data)
	}
	if !reflect.DeepEqual(org.Program, rst.Program) {
		t.Errorf("Program:got %v, expected %v\n", rst.Program, org.Program)
	}
}

func TestFSTOperationString(t *testing.T) {

	ps := []struct {
		op   Operation
		name string
	}{
		{0, "UNDEF0"},
		{1, "ACCEPT"},
		{2, "ACCEPTB"},
		{3, "MATCH"},
		{4, "MATCHB"},
		{5, "OUTPUT"},
		{6, "OUTPUTB"},
		{7, "UNDEF7"},
		{8, "NA[8]"},
		{9, "NA[9]"},
	}

	for _, p := range ps {
		if p.op.String() != p.name {
			t.Errorf("got %v, expected %v", p.op.String(), p.name)
		}
	}
}

func TestFSTStress(t *testing.T) {
	fp, err := os.Open("../testdata/words.txt")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	defer fp.Close()

	var ps PairSlice
	s := bufio.NewScanner(fp)
	for i := 0; s.Scan(); i++ {
		p := Pair{In: s.Text(), Out: int32(i)}
		ps = append(ps, p)
	}
	if e := s.Err(); e != nil {
		t.Fatalf("unexpected error, %v", e)
	}
	m := BuildMAST(ps)
	fst, err := m.BuildFST()
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	for _, p := range ps {
		ids := fst.Search(p.In)
		if !func(s []int32, x int32) bool {
			for i := range s {
				if x == s[i] {
					return true
				}
			}
			return false
		}(ids, p.Out) {
			t.Errorf("input:%v, got %v, but not input %v", p.In, ids, p.Out)
		}
	}
}
