package string

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func (t *FST) runTester(input string) (cs []Configuration, accept bool) {
	t.Run(input, func(snapshot Configuration) {
		cs = append(cs, snapshot)
		accept = snapshot.Head == len(input)
	})
	return cs, accept
}

func (t FST) searchTester(input string) []string {
	snap, acc := t.runTester(input)
	if !acc || len(snap) == 0 {
		return nil
	}
	c := snap[len(snap)-1]
	return c.Outputs
}

func (t FST) prefixSearchTester(input string) (length int, output []string) {
	snap, _ := t.runTester(input)
	if len(snap) == 0 {
		return -1, nil
	}
	c := snap[len(snap)-1]
	return c.Head, c.Outputs
}

func (t FST) commonPrefixSearchTester(input string) (lens []int, outputs [][]string) {
	snap, _ := t.runTester(input)
	if len(snap) == 0 {
		return lens, outputs
	}
	for _, c := range snap {
		lens = append(lens, c.Head)
		outputs = append(outputs, c.Outputs)
	}
	return lens, outputs

}

func TestFSTRun01(t *testing.T) {
	input := PairSlice{
		{In: "feb", Out: "28"},
		{In: "feb", Out: "29"},
		{In: "feb", Out: "30"},
		{In: "dec", Out: "31"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

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
				Outputs: []string{"28", "29", "30"},
			},
		}
		got, ok := fst.runTester("feb")
		if !ok {
			t.Errorf("input:feb, config:%v, Accept:%v", got, ok)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})
}

func TestFSTRun02(t *testing.T) {
	input := PairSlice{
		{In: "feb", Out: "28"},
		{In: "feb", Out: "29"},
		{In: "feb", Out: "30"},
		{In: "dec", Out: "31"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

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
				Outputs: []string{"31"},
			},
		}
		got, ok := fst.runTester("dec")
		if !ok {
			t.Errorf("input:dec, config:%v, Accept:%v", got, ok)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})
}

func TestFSTRun03(t *testing.T) {
	input := PairSlice{
		{In: "feb", Out: "hello,world"},
		{In: "february", Out: "hell world"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	t.Run("machine code", func(t *testing.T) {
		expected := `  0 OUTPUTB	66(f) 1
  1 [hell]
  2 MATCHB	65(e) 1
  3 MATCHB	62(b) 1
  4 ACCEPT	1 0
  5 [2]
  6 [1] [o,world]
  7 OUTPUTB	72(r) 1
  8 [ world]
  9 MATCHB	75(u) 1
 10 MATCHB	61(a) 1
 11 MATCHB	72(r) 1
 12 MATCHB	79(y) 1
 13 ACCEPTB	0 0
`
		if got := fmt.Sprint(fst); got != expected {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("running", func(t *testing.T) {
		expected := []Configuration{
			{
				PC:      4,
				Head:    3,
				Outputs: []string{"hello,world"},
			},
			{
				PC:      13,
				Head:    8,
				Outputs: []string{"hell world"},
			},
		}
		got, ok := fst.runTester("february")
		if !ok {
			t.Errorf("input:dec, config:%v, Accept:%v", got, ok)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})
}

func TestFstSearch01(t *testing.T) {
	input := PairSlice{
		{In: "1a22xss", Out: "111"},
		{In: "1a22", Out: "111"},
		{In: "1b22yss", Out: "222"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	for _, p := range input {
		outs := fst.searchTester(p.In)
		if !reflect.DeepEqual(outs, []string{p.Out}) {
			t.Errorf("input %v, got %v, expected %v", p.In, outs, []string{p.Out})
		}
	}
}

func TestFstVMSearch02(t *testing.T) {
	input := PairSlice{
		{In: "hell", Out: "666"},
		{In: "hello", Out: "111"},
		{In: "goodbye", Out: "222"},
		{In: "goodbye", Out: "333"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	testdata := []struct {
		input    string
		expected []string
	}{
		{input: "hell", expected: []string{"666"}},
		{input: "hello", expected: []string{"111"}},
		{input: "goodbye", expected: []string{"222", "333"}},
	}
	for _, d := range testdata {
		outs := fst.searchTester(d.input)
		if !reflect.DeepEqual(outs, d.expected) {
			t.Errorf("input %v, got %v, expected %v", d.input, outs, d.expected)
		}
	}
}

func TestFstSearch03(t *testing.T) {
	input := PairSlice{
		{In: "hell", Out: ""},
		{In: "hello", Out: ""},
		{In: "goodbye", Out: ""},
		{In: "goodbye", Out: ""},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	testdata := []struct {
		input    string
		expected []string
	}{
		{input: "hell", expected: []string{""}},
		{input: "hello", expected: []string{""}},
		{input: "goodbye", expected: []string{""}},
	}
	for _, d := range testdata {
		outs := fst.searchTester(d.input)
		if !reflect.DeepEqual(outs, d.expected) {
			t.Errorf("input %v, got %v, expected %v", d.input, outs, d.expected)
		}
	}
}

func TestFstSearch04(t *testing.T) {
	input := PairSlice{
		{In: "1a22", Out: "aloha"},
		{In: "1a22xss", Out: "world"},
		{In: "1a22xss", Out: "goodby"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	in := "1a22xss"
	expected := []string{"goodby", "world"}
	got := fst.searchTester(in)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("input: %v, got %v, expected %v", in, got, expected)
	}
}

func TestFstSearch05(t *testing.T) {
	input := PairSlice{
		{"1a22", "goodbye"},
		{"1a22xss", "goodbye"},
		{"1a22xss", "good"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	in := "1a22xss"
	expected := []string{"good", "goodbye"}
	got := fst.searchTester(in)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("input: %v, got %v, expected %v", in, got, expected)
	}
}

func TestFstSearch06(t *testing.T) {
	input := PairSlice{
		{In: "こんにちは", Out: "hello"},
		{In: "世界", Out: "world"},
		{In: "すもももももも", Out: "peach"},
		{In: "すもも", Out: "peach"},
		{In: "すもも", Out: "もも"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	testdata := []struct {
		in       string
		expected []string
	}{
		{in: "すもも", expected: []string{"peach", "もも"}},
		{in: "こんにちわ"},
		{in: "こんにちは", expected: []string{"hello"}},
		{in: "世界", expected: []string{"world"}},
		{in: "すもももももも", expected: []string{"peach"}},
		{in: "すも", expected: nil},
		{in: "すもう", expected: nil},
	}

	for _, d := range testdata {
		got := fst.searchTester(d.in)
		if !reflect.DeepEqual(got, d.expected) {
			t.Errorf("input:%v, got %v, expected %v", d.in, got, d.expected)
		}
	}
}

func TestFstPrefixSearch01(t *testing.T) {
	input := PairSlice{
		{"こんにちは", "hello"},
		{"世界", "world"},
		{"すもももももも", "peach"},
		{"すもも", "peach"},
		{"すもも", "もも"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	testdata := []struct {
		in  string
		pos int
		out []string
	}{
		{in: "すもも", pos: 9, out: []string{"peach", "もも"}},
		{in: "こんにちわ", pos: -1},
		{in: "こんにちは", pos: 15, out: []string{"hello"}},
		{in: "世界", pos: 6, out: []string{"world"}},
		{in: "すもももももも", pos: 21, out: []string{"peach"}},
		{in: "すも", pos: -1, out: nil},
		{in: "すもう", pos: -1, out: nil},
		{in: "すもももももももものうち", pos: 21, out: []string{"peach"}},
	}

	for _, d := range testdata {
		pos, outs := fst.prefixSearchTester(d.in)
		sort.Strings(outs)
		sort.Strings(d.out)
		if pos != d.pos || !reflect.DeepEqual(outs, d.out) {
			t.Errorf("input:%v, got %v %v, expected %v %v", d.in, pos, outs, d.pos, d.out)
		}
	}
}

func TestFstVMCommonPrefixSearch01(t *testing.T) {
	input := PairSlice{
		{"こんにちは", "hello"},
		{"世界", "world"},
		{"すもももももも", "peach"},
		{"すもも", "peach"},
		{"すもも", "もも"},
	}
	m := BuildMAST(input)
	fst, err := m.BuildFST()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	testdata := []struct {
		in   string
		lens []int
		outs [][]string
	}{
		{"すもも", []int{9}, [][]string{{"peach", "もも"}}},
		{"こんにちわ", nil, nil},
		{"こんにちは", []int{15}, [][]string{{"hello"}}},
		{"世界", []int{6}, [][]string{{"world"}}},
		{"すもももももも", []int{9, 21}, [][]string{{"peach", "もも"}, {"peach"}}},
		{"すも", nil, nil},
		{"すもう", nil, nil},
		{"すもももももももものうち", []int{9, 21}, [][]string{{"peach", "もも"}, {"peach"}}},
	}

	for _, d := range testdata {
		lens, outs := fst.commonPrefixSearchTester(d.in)
		if !reflect.DeepEqual(lens, d.lens) || len(outs) != len(d.outs) {
			t.Errorf("input:%v, got lens:%v outs:%v, expected lens:%v outs:%v",
				d.in, lens, outs, d.lens, d.outs)
			continue
		}
		for i := range outs {
			got := outs[i]
			expected := d.outs[i]
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("input:%v, got lens:%v outs:%v, expected lens:%v outs:%v",
					d.in, lens, outs, d.lens, d.outs)
			}
		}
	}
}

func TestSaveAndLoad01(t *testing.T) {
	input := PairSlice{
		{"こんにちは", "hello"},
		{"世界", "world"},
		{"すもももももも", "peach"},
		{"すもも", "peach"},
		{"すもも", "もも"},
	}
	m := BuildMAST(input)
	src, e := m.BuildFST()
	if e != nil {
		t.Errorf("unexpected error %v", e)
	}
	var b bytes.Buffer
	if _, err := src.WriteTo(&b); err != nil {
		t.Errorf("unexpected error %v", e)
	}

	dst, err := Read(&b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(src.Program, dst.Program) {
		t.Errorf("saved program:\n%v\nloaded program:\n%v", src, dst)
	}
	if !reflect.DeepEqual(src.Data, dst.Data) {
		t.Errorf("saved data:\n%v\nloaded data:\n%v", src.Data, dst.Data)
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
		p := Pair{In: s.Text(), Out: strconv.Itoa(i)}
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
		ids := fst.searchTester(p.In)
		if !func(s []string, x string) bool {
			for i := range s {
				if x == s[i] {
					return true
				}
			}
			return false
		}(ids, p.Out) {
			t.Errorf("input:%v, got %v, but not in %v", p.In, ids, p.Out)
		}
	}
}
