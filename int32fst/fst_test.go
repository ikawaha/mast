package int32fst

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/ikawaha/dartsclone"
)

func (t *FST) runTester(input string) (cs []Configuration, accept bool) {
	t.Run(input, func(snapshot Configuration) {
		cs = append(cs, snapshot)
		accept = snapshot.Head == len(input)
	})
	return cs, accept
}

func TestFSTRun01(t *testing.T) {
	input := PairSlice{
		{In: "feb", Out: 28},
		{In: "feb", Out: 29},
		{In: "feb", Out: 30},
		{In: "dec", Out: 31},
	}
	fst, err := New(input)
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
				Outputs: []int32{28, 29, 30},
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
	inp := PairSlice{
		{In: "feb", Out: 28},
		{In: "feb", Out: 29},
		{In: "feb", Out: 30},
		{In: "dec", Out: 31},
	}
	m := BuildMAST(inp)
	fst, err := BuildFST(m)
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
				Outputs: []int32{31},
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
	inp := PairSlice{
		{In: "feb", Out: 0},
		{In: "february", Out: maxUint16 + 1},
	}
	m := BuildMAST(inp)
	fst, err := BuildFST(m)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	t.Run("machine code", func(t *testing.T) {
		expected := `  0 MATCHB	66(f) 1
  1 MATCHB	65(e) 1
  2 MATCHB	62(b) 1
  3 ACCEPT	0 0
  4 OUTPUTB	72(r) 1
  5 [65536]
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
				Outputs: []int32{65536},
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

func TestFSTSearch01(t *testing.T) {
	input := PairSlice{
		{In: "1a22xss", Out: 111},
		{In: "1a22", Out: 111},
		{In: "1b22yss", Out: 222},
	}
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	for _, p := range input {
		outs := fst.Search(p.In)
		if !reflect.DeepEqual(outs, []int32{p.Out}) {
			t.Errorf("input %v, got %v, expected %v", p.In, outs, []int32{p.Out})
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
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
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
			t.Errorf("input %v, got %v, expected %v", d.input, outs, d.expected)
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
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
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
			t.Errorf("input %v, got %v, expected %v", d.input, outs, d.expected)
		}
	}
}

func TestFSTSearch04(t *testing.T) {
	input := PairSlice{
		{In: "こんにちは", Out: 111},
		{In: "世界", Out: 222},
		{In: "すもももももも", Out: 333},
		{In: "すもも", Out: 333},
		{In: "すもも", Out: 444},
	}
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
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
			t.Errorf("input:%v, got %v, expected %v", d.input, outs, d.expected)
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
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
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
			t.Errorf("input:%v, got %v, expected %v", d.input, outs, d.outputs)
		}
		if pos != d.pos {
			t.Errorf("input:%v, got %v, expected %v", d.input, pos, d.pos)
		}
	}
}

func TestFSTCommonPrefixSearch01(t *testing.T) {
	input := PairSlice{
		{In: "こんにちは", Out: 111},
		{In: "世界", Out: 222},
		{In: "すもももももも", Out: 333},
		{In: "すもも", Out: 333},
		{In: "すもも", Out: 444},
	}
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
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
			t.Errorf("input:%v, got %v %v, expected %v %v", d.input, lens, outs, d.lens, d.outputs)
		}
		for i := range lens {
			sort.Sort(int32Slice(outs[i]))
			sort.Sort(int32Slice(d.outputs[i]))
			if !reflect.DeepEqual(outs[i], d.outputs[i]) {
				t.Errorf("input:%v, got %v %v, expected %v %v", d.input, lens, outs, d.lens, d.outputs)
			}
		}
	}
}

func TestFSTSaveAndLoad01(t *testing.T) {
	var input = PairSlice{
		{In: "feb", Out: 28},
		{In: "feb", Out: 29},
		{In: "apr", Out: 30},
		{In: "jan", Out: 31},
		{In: "jun", Out: 30},
		{In: "jul", Out: 31},
		{In: "dec", Out: 31},
	}
	m := BuildMAST(input)
	src, err := BuildFST(m)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var b bytes.Buffer
	n, err := src.WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	if n != int64(b.Len()) {
		t.Errorf("write len: got %v, expected %v", n, b.Len())
	}

	dst, err := Read(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(src.Data, dst.Data) {
		t.Errorf("Data:got %v, expected %v", dst.Data, src.Data)
	}
	if !reflect.DeepEqual(src.Program, dst.Program) {
		t.Errorf("Program:got %v, expected %v", dst.Program, src.Program)
	}
}

func TestFSTOperationString(t *testing.T) {
	testdata := []struct {
		op   Operation
		name string
	}{
		{op: 0, name: "UNDEF0"},
		{op: 1, name: "ACCEPT"},
		{op: 2, name: "ACCEPTB"},
		{op: 3, name: "MATCH"},
		{op: 4, name: "MATCHB"},
		{op: 5, name: "OUTPUT"},
		{op: 6, name: "OUTPUTB"},
		{op: 7, name: "UNDEF7"},
		{op: 8, name: "NA[8]"},
		{op: 9, name: "NA[9]"},
	}

	for _, d := range testdata {
		if d.op.String() != d.name {
			t.Errorf("got %v, expected %v", d.op.String(), d.name)
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
		p := Pair{In: s.Text(), Out: maxUint16 + int32(i)}
		ps = append(ps, p)
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	m := BuildMAST(ps)
	fst, err := BuildFST(m)
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

func TestSearch(t *testing.T) {
	input := PairSlice{
		{In: "こんにちは", Out: 111},
		{In: "世界", Out: 222},
		{In: "すもももももも", Out: 333},
		{In: "すもも", Out: 333},
	}
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	for _, v := range input {
		expected := []int32{v.Out}
		if got := fst.Search(v.In); !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v, %v", got, expected, v.In)
		}
	}
	// expected not to be found.
	if got := fst.Search("すももも"); got != nil {
		t.Errorf("got %v, expected nil", got)
	}
}

func TestPrefixSearch(t *testing.T) {
	var input = PairSlice{
		{In: "東京", Out: 1},
		{In: "東京チョコレート", Out: 2},
		{In: "東京チョコレートMIX", Out: 3},
		{In: "hello", Out: 4},
		{In: "goodbye", Out: 5},
		{In: "good", Out: 6},
		{In: "go", Out: 7},
		{In: "gopher", Out: 8},
	}
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Run("東京チョコレートMIX", func(t *testing.T) {
		length, outs := fst.PrefixSearch("東京チョコレートMIX!!!")
		if expected := len("東京チョコレートMIX"); length != expected {
			t.Errorf("got %v, expected %v", length, expected)
		}
		if expected := []int32{3}; !reflect.DeepEqual(outs, expected) {
			t.Errorf("got %v, expected %v", outs, expected)
		}
	})

	t.Run("good-by", func(t *testing.T) {
		length, outs := fst.PrefixSearch("good-by")
		if expected := len("good"); length != expected {
			t.Errorf("got %v, expected %v", length, expected)
		}
		if expected := []int32{6}; !reflect.DeepEqual(outs, expected) {
			t.Errorf("got %v, expected %v", outs, expected)
		}
	})

	t.Run("not found", func(t *testing.T) {
		length, outs := fst.PrefixSearch("aloha")
		if expected := -1; length != expected {
			t.Errorf("got %v, expected %v", length, expected)
		}
		if outs != nil {
			t.Errorf("got %v, expected nil", outs)
		}
	})
}

func TestCommonPrefixSearch(t *testing.T) {
	input := PairSlice{
		{In: "東京", Out: 1},
		{In: "東京チョコレート", Out: 2},
		{In: "東京チョコレートMIX", Out: 3},
		{In: "hello", Out: 4},
		{In: "goodbye", Out: 5},
		{In: "good", Out: 6},
		{In: "go", Out: 7},
		{In: "go", Out: 77},
		{In: "gopher", Out: 8},
	}
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Run("東京チョコレートMIX", func(t *testing.T) {
		lens, outs := fst.CommonPrefixSearch("東京チョコレートMIX!!!")
		expectedLens := []int{len("東京"), len("東京チョコレート"), len("東京チョコレートMIX")}
		if !reflect.DeepEqual(lens, expectedLens) {
			t.Errorf("got %v, expected %v", lens, expectedLens)
		}
		expectedOuts := [][]int32{{1}, {2}, {3}}
		if !reflect.DeepEqual(outs, expectedOuts) {
			t.Errorf("got %v, expected %v", outs, expectedOuts)
		}
	})

	t.Run("good-by", func(t *testing.T) {
		lens, outs := fst.CommonPrefixSearch("good-by")
		expectedLens := []int{len("go"), len("good")}
		if !reflect.DeepEqual(lens, expectedLens) {
			t.Errorf("got %v, expected %v", lens, expectedLens)
		}
		expectedOuts := [][]int32{{7, 77}, {6}}
		if !reflect.DeepEqual(outs, expectedOuts) {
			t.Errorf("got %v, expected %v", outs, expectedOuts)
		}
	})

	t.Run("not found", func(t *testing.T) {
		lens, outs := fst.CommonPrefixSearch("aloha")
		if lens != nil {
			t.Errorf("got %v, expected nil", lens)
		}
		if outs != nil {
			t.Errorf("got %v, expected nil", outs)
		}
	})
}

func TestCommonPrefixSearchCallback(t *testing.T) {
	var input = PairSlice{
		{In: "東京", Out: 1},
		{In: "東京チョコレート", Out: 2},
		{In: "東京チョコレートMIX", Out: 3},
		{In: "hello", Out: 4},
		{In: "goodbye", Out: 5},
		{In: "good", Out: 6},
		{In: "go", Out: 7},
		{In: "go", Out: 77},
		{In: "gopher", Out: 8},
	}
	fst, err := New(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Run("東京チョコレートMIX", func(t *testing.T) {
		var (
			lens []int
			outs [][]int32
		)
		fst.CommonPrefixSearchCallback("東京チョコレートMIX!!!", func(length int, outputs []int32) {
			lens = append(lens, length)
			outs = append(outs, outputs)
		})
		expectedLens := []int{len("東京"), len("東京チョコレート"), len("東京チョコレートMIX")}
		if !reflect.DeepEqual(lens, expectedLens) {
			t.Errorf("got %v, expected %v", lens, expectedLens)
		}
		expectedOuts := [][]int32{{1}, {2}, {3}}
		if !reflect.DeepEqual(outs, expectedOuts) {
			t.Errorf("got %v, expected %v", outs, expectedOuts)
		}
	})

	t.Run("good-by", func(t *testing.T) {
		var (
			lens []int
			outs [][]int32
		)
		fst.CommonPrefixSearchCallback("good-by", func(length int, outputs []int32) {
			lens = append(lens, length)
			outs = append(outs, outputs)
		})
		expectedLens := []int{len("go"), len("good")}
		if !reflect.DeepEqual(lens, expectedLens) {
			t.Errorf("got %v, expected %v", lens, expectedLens)
		}
		expectedOuts := [][]int32{{7, 77}, {6}}
		if !reflect.DeepEqual(outs, expectedOuts) {
			t.Errorf("got %v, expected %v", outs, expectedOuts)
		}
	})

	t.Run("not found", func(t *testing.T) {
		fst.CommonPrefixSearchCallback("aloha", func(length int, outputs []int32) {
			// expects not to call
			t.Errorf("unecpected call, length %v, outputs %v", length, outputs)
		})
	})
}

func BenchmarkSearch(b *testing.B) {
	fp, err := os.Open("../testdata/ipadic.txt")
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	defer fp.Close()
	var input PairSlice
	s := bufio.NewScanner(fp)
	for i := 0; s.Scan(); i++ {
		p := Pair{In: s.Text(), Out: int32(i)}
		input = append(input, p)
	}
	if err := s.Err(); err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	fst, err := New(input)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}

	xxx, err := os.Create("fstidpadicxxx")
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	defer xxx.Close()
	fst.WriteTo(xxx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range input {
			fst.Search(v.In)
		}
	}
}

func BenchmarkSearchDA(b *testing.B) {
	fp, err := os.Open("../testdata/words_uniq.txt")
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	defer fp.Close()
	var input []string
	s := bufio.NewScanner(fp)
	for i := 0; s.Scan(); i++ {
		input = append(input, s.Text())
	}
	if err := s.Err(); err != nil {
		b.Fatalf("unexpected error, %v", err)
	}

	trie, err := dartsclone.BuildTRIE(input, nil, nil)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}

	builder := dartsclone.NewBuilder(nil)
	builder.Build(input, nil)
	xxx, err := os.Create("daxxx")
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	builder.WriteTo(xxx)
	xxx.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range input {
			trie.ExactMatchSearch(v)
		}
	}
}
