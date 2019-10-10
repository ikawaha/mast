package ss

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestFSTRun01(t *testing.T) {
	inp := PairSlice{
		{"feb", "28"},
		{"feb", "29"},
		{"feb", "30"},
		{"dec", "31"},
	}
	m := BuildMAST(inp)
	m.Dot(os.Stdout)

	fst, _ := m.BuildFST()
	fmt.Println(fst)

	config, accept := fst.Run("feb")
	if !accept {
		t.Errorf("input:feb, config:%v, Accept:%v", config, accept)
	}
	fmt.Println(config)

}

func TestFSTRun02(t *testing.T) {
	inp := PairSlice{
		{"feb", "28"},
		{"feb", "29"},
		{"feb", "30"},
		{"dec", "31"},
	}
	m := BuildMAST(inp)
	m.Dot(os.Stdout)

	fst, _ := m.BuildFST()
	fmt.Println(fst)

	config, ok := fst.Run("dec")
	if !ok {
		t.Errorf("input:feb, config:%v, Accept:%v", config, ok)
	}
	fmt.Println(config)

}

func TestFSTRun03(t *testing.T) {
	inp := PairSlice{
		{"feb", "0"},
		{"february", "1"},
	}
	m := BuildMAST(inp)
	m.Dot(os.Stdout)

	fst, _ := m.BuildFST()
	fmt.Println(fst)

	input := "february"
	config, ok := fst.Run(input)
	if !ok {
		t.Errorf("input:%v, config:%+v, Accept:%v", input, config, ok)
	}
	fmt.Println(config)
}

func TestFstSearch01(t *testing.T) {
	inp := PairSlice{
		{"1a22xss", "world"},
		{"1b22yss", "goodby"},
	}
	fst, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(fst)
	for _, p := range inp {
		outs := fst.Search(p.In)
		if !reflect.DeepEqual(outs, []string{p.Out}) {
			t.Errorf("input: %v, got %v, expected %v\n", p.In, outs, []string{p.Out})
		}
	}
}

func TestFstVMSearch02(t *testing.T) {
	inp := PairSlice{
		{"1a22", "aloha"},
		{"1a22xss", "world"},
		{"1a22yss", "goodby"},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)
	for _, p := range inp {
		outs := vm.Search(p.In)
		if !reflect.DeepEqual(outs, []string{p.Out}) {
			t.Errorf("input: %v, got %v, expected %v\n", p.In, outs, []string{p.Out})
		}
	}
}

func TestFstVMSearch03(t *testing.T) {
	inp := PairSlice{
		{"1a22", "aloha"},
		{"1a22xss", "world"},
		{"1a22xss", "goodby"},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22"
	exp := []string{"aloha"}
	outs := vm.Search(in)
	if !reflect.DeepEqual(outs, []string{"aloha"}) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFstVMSearch04(t *testing.T) {
	inp := PairSlice{
		{"1a22", "aloha"},
		{"1a22xss", "world"},
		{"1a22xss", "goodby"},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22xss"
	exp := []string{"goodby", "world"}
	outs := vm.Search(in)
	sort.Strings(outs)
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFstVMSearch05(t *testing.T) {
	inp := PairSlice{
		{"1a22", ""},
		{"1a22xss", ""},
		{"1a22xss", ""},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22xss"
	exp := []string{""}
	outs := vm.Search(in)
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFstVMSearch06(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", "hello"},
		{"世界", "world"},
		{"すもももももも", "peach"},
		{"すもも", "peach"},
		{"すもも", "もも"},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	cr := []struct {
		in  string
		out []string
	}{
		{"すもも", []string{"peach", "もも"}},
		{"こんにちわ", nil},
		{"こんにちは", []string{"hello"}},
		{"世界", []string{"world"}},
		{"すもももももも", []string{"peach"}},
		{"すも", nil},
		{"すもう", nil},
	}

	for _, pair := range cr {
		outs := vm.Search(pair.in)
		sort.Strings(outs)
		sort.Strings(pair.out)
		if !reflect.DeepEqual(outs, pair.out) {
			t.Errorf("input:%v, got %v, expected %v\n", pair.in, outs, pair.out)
		}
	}
}

func TestFstVMPrefixSearch01(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", "hello"},
		{"世界", "world"},
		{"すもももももも", "peach"},
		{"すもも", "peach"},
		{"すもも", "もも"},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	crs := []struct {
		in  string
		pos int
		out []string
	}{
		{"すもも", 9, []string{"peach", "もも"}},
		{"こんにちわ", -1, nil},
		{"こんにちは", 15, []string{"hello"}},
		{"世界", 6, []string{"world"}},
		{"すもももももも", 21, []string{"peach"}},
		{"すも", -1, nil},
		{"すもう", -1, nil},
		{"すもももももももものうち", 21, []string{"peach"}},
	}

	for _, cr := range crs {
		pos, outs := vm.PrefixSearch(cr.in)
		sort.Strings(outs)
		sort.Strings(cr.out)
		if pos != cr.pos || !reflect.DeepEqual(outs, cr.out) {
			t.Errorf("input:%v, got %v %v, expected %v %v\n", cr.in, pos, outs, cr.pos, cr.out)
		}
	}
}

func TestFstVMCommonPrefixSearch01(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", "hello"},
		{"世界", "world"},
		{"すもももももも", "peach"},
		{"すもも", "peach"},
		{"すもも", "もも"},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	crs := []struct {
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

	for _, cr := range crs {
		lens, outs := vm.CommonPrefixSearch(cr.in)
		if !reflect.DeepEqual(lens, cr.lens) || len(outs) != len(cr.outs) {
			t.Errorf("input:%v, got lens:%v outs:%v, expected lens:%v outs:%v",
				cr.in, lens, outs, cr.lens, cr.outs)
			continue
		}
		for i := range outs {
			o := outs[i]
			e := cr.outs[i]
			sort.Strings(o)
			sort.Strings(e)
			if !reflect.DeepEqual(o, e) {
				t.Errorf("input:%v, got lens:%v outs:%v, expected lens:%v outs:%v",
					cr.in, lens, outs, cr.lens, cr.outs)
			}
		}
	}
}

func TestSaveAndLoad01(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", "hello"},
		{"世界", "world"},
		{"すもももももも", "peach"},
		{"すもも", "peach"},
		{"すもも", "もも"},
	}
	src, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error %v\n", e)
	}
	var b bytes.Buffer
	if _, err := src.WriteTo(&b); err != nil {
		t.Errorf("unexpected error %v\n", e)
	}

	dst, err := Read(&b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(src.Program, dst.Program) {
		t.Errorf("saved program:\n%v\nloaded program:\n%v\n", src, dst)
	}
	if !reflect.DeepEqual(src.Data, dst.Data) {
		t.Errorf("saved data:\n%v\nloaded data:\n%v\n", src.Data, dst.Data)
	}
}

/*
func TestFSTSearch02(t *testing.T) {
	inp := PairSlice{
		{"1a22", 111},
		{"1a22xss", 222},
		{"1a22yss", 333},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)
	for _, p := range inp {
		outs := vm.Search(p.In)
		if !reflect.DeepEqual(outs, []int32{p.Out}) {
			t.Errorf("input: %v, got %v, expected %v\n", p.In, outs, []int32{p.Out})
		}
	}
}

func TestFSTSearch03(t *testing.T) {
	inp := PairSlice{
		{"1a22", 111},
		{"1a22xss", 222},
		{"1a22xss", 333},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22"
	exp := []int32{111}
	outs := vm.Search(in)
	if !reflect.DeepEqual(outs, []int32{111}) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFSTSearch04(t *testing.T) {
	inp := PairSlice{
		{"1a22", 111},
		{"1a22xss", 222},
		{"1a22xss", 333},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22xss"
	exp := []int32{222, 333}
	outs := vm.Search(in)
	sort.Sort(int32Slice(outs))
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFSTSearch05(t *testing.T) {
	inp := PairSlice{
		{"1a22", 0},
		{"1a22xss", 0},
		{"1a22xss", 0},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22xss"
	exp := []int32{0}
	outs := vm.Search(in)
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFSTSearch06(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", 111},
		{"世界", 222},
		{"すもももももも", 333},
		{"すもも", 333},
		{"すもも", 444},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	cr := []struct {
		in  string
		out []int32
	}{
		{"すもも", []int32{333, 444}},
		{"こんにちわ", nil},
		{"こんにちは", []int32{111}},
		{"世界", []int32{222}},
		{"すもももももも", []int32{333}},
		{"すももももももも", nil},
		{"すも", nil},
		{"すもう", nil},
	}

	for _, pair := range cr {
		outs := vm.Search(pair.in)
		sort.Sort(int32Slice(outs))
		sort.Sort(int32Slice(pair.out))
		if !reflect.DeepEqual(outs, pair.out) {
			t.Errorf("input:%v, got %v, expected %v\n", pair.in, outs, pair.out)
		}
	}
}

func TestFSTPrefixSearch01(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", 111},
		{"世界", 222},
		{"すもももももも", 333},
		{"すもも", 333},
		{"すもも", 444},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	crs := []struct {
		in  string
		pos int
		out []int32
	}{
		{"すもも", 9, []int32{333, 444}},
		{"こんにちわ", -1, nil},
		{"こんにちは", 15, []int32{111}},
		{"世界", 6, []int32{222}},
		{"すもももももも", 21, []int32{333}},
		{"すもももももももものうち", 21, []int32{333}},
		{"すも", -1, nil},
		{"すもう", -1, nil},
	}

	for _, cr := range crs {
		pos, outs := vm.PrefixSearch(cr.in)
		sort.Sort(int32Slice(outs))
		sort.Sort(int32Slice(cr.out))
		if !reflect.DeepEqual(outs, cr.out) {
			t.Errorf("input:%v, got %v, expected %v\n", cr.in, outs, cr.out)
		}
		if pos != cr.pos {
			t.Errorf("input:%v, got %v, expected %v\n", cr.in, pos, cr.pos)
		}
	}
}

func TestFSTCommonPrefixSearch01(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", 111},
		{"世界", 222},
		{"すもももももも", 333},
		{"すもも", 333},
		{"すもも", 444},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	crs := []struct {
		in   string
		lens []int
		outs [][]int32
	}{
		{"すもも", []int{9}, [][]int32{{333, 444}}},
		{"こんにちわ", nil, nil},
		{"こんにちは", []int{15}, [][]int32{{111}}},
		{"世界", []int{6}, [][]int32{{222}}},
		{"すもももももも", []int{9, 21}, [][]int32{{333, 444}, {333}}},
		{"すもももももももものうち", []int{9, 21}, [][]int32{{333, 444}, {333}}},
		{"すも", nil, nil},
		{"すもう", nil, nil},
	}

	for _, cr := range crs {
		lens, outs := vm.CommonPrefixSearch(cr.in)
		if !reflect.DeepEqual(lens, cr.lens) {
			t.Errorf("input:%v, got %v %v, expected %v %v\n", cr.in, lens, outs, cr.lens, cr.outs)
		}
		for i := range lens {
			sort.Sort(int32Slice(outs[i]))
			sort.Sort(int32Slice(cr.outs[i]))
			if !reflect.DeepEqual(outs[i], cr.outs[i]) {
				t.Errorf("input:%v, got %v %v, expected %v %v\n", cr.in, lens, outs, cr.lens, cr.outs)
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

	org, e := Build(inp)
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
	fp, err := os.Open("./testdata/words.txt")
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
			t.Errorf("input:%v, got %v, but not in %v", p.In, ids, p.Out)
		}
	}
}
*/
