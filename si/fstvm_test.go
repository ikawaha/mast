package si

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestFstVMSearch01(t *testing.T) {
	inp := PairSlice{
		{"1a22xss", 111},
		{"1b22yss", 222},
	}
	vm, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)
	for _, p := range inp {
		outs := vm.Search(p.In)
		if !reflect.DeepEqual(outs, []int{p.Out}) {
			t.Errorf("input: %v, got %v, expected %v\n", p.In, outs, []int{p.Out})
		}
	}
}

func TestFstVMSearch02(t *testing.T) {
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
		if !reflect.DeepEqual(outs, []int{p.Out}) {
			t.Errorf("input: %v, got %v, expected %v\n", p.In, outs, []int{p.Out})
		}
	}
}

func TestFstVMSearch03(t *testing.T) {
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
	exp := []int{111}
	outs := vm.Search(in)
	if !reflect.DeepEqual(outs, []int{111}) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFstVMSearch04(t *testing.T) {
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
	exp := []int{222, 333}
	outs := vm.Search(in)
	sort.Ints(outs)
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFstVMSearch05(t *testing.T) {
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
	exp := []int{0}
	outs := vm.Search(in)
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFstVMSearch06(t *testing.T) {
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
		out []int
	}{
		{"すもも", []int{333, 444}},
		{"こんにちわ", nil},
		{"こんにちは", []int{111}},
		{"世界", []int{222}},
		{"すもももももも", []int{333}},
		{"すももももももも", nil},
		{"すも", nil},
		{"すもう", nil},
	}

	for _, pair := range cr {
		outs := vm.Search(pair.in)
		sort.Ints(outs)
		sort.Ints(pair.out)
		if !reflect.DeepEqual(outs, pair.out) {
			t.Errorf("input:%v, got %v, expected %v\n", pair.in, outs, pair.out)
		}
	}
}

func TestFstVMPrefixSearch06(t *testing.T) {
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
		out []int
	}{
		{"すもも", 9, []int{333, 444}},
		{"こんにちわ", -1, nil},
		{"こんにちは", 15, []int{111}},
		{"世界", 6, []int{222}},
		{"すもももももも", 21, []int{333}},
		{"すもももももももものうち", 21, []int{333}},
		{"すも", -1, nil},
		{"すもう", -1, nil},
	}

	for _, cr := range crs {
		pos, outs := vm.PrefixSearch(cr.in)
		sort.Ints(outs)
		sort.Ints(cr.out)
		if !reflect.DeepEqual(outs, cr.out) {
			t.Errorf("input:%v, got %v, expected %v\n", cr.in, outs, cr.out)
		}
		if pos != cr.pos {
			t.Errorf("input:%v, got %v, expected %v\n", cr.in, pos, cr.pos)
		}
	}
}

func TestFstVMCommonPrefixSearch06(t *testing.T) {
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
		outs [][]int
	}{
		{"すもも", []int{9}, [][]int{{333, 444}}},
		{"こんにちわ", nil, nil},
		{"こんにちは", []int{15}, [][]int{{111}}},
		{"世界", []int{6}, [][]int{{222}}},
		{"すもももももも", []int{9, 21}, [][]int{{333, 444}, {333}}},
		{"すもももももももものうち", []int{9, 21}, [][]int{{333, 444}, {333}}},
		{"すも", nil, nil},
		{"すもう", nil, nil},
	}

	for _, cr := range crs {
		lens, outs := vm.CommonPrefixSearch(cr.in)
		if !reflect.DeepEqual(lens, cr.lens) {
			t.Errorf("input:%v, got %v %v, expected %v %v\n", cr.in, lens, outs, cr.lens, cr.outs)
		}
		for i := range lens {
			sort.Ints(outs[i])
			sort.Ints(cr.outs[i])
			if !reflect.DeepEqual(outs[i], cr.outs[i]) {
				t.Errorf("input:%v, got %v %v, expected %v %v\n", cr.in, lens, outs, cr.lens, cr.outs)
			}
		}
	}
}

func TestToInt01(t *testing.T) {
	cr := []struct {
		in  []byte
		out int
	}{
		{[]byte{7}, 7},
		{[]byte{255}, 255},
		{[]byte{0, 1}, 1},
		{[]byte{1, 0}, 256},
		{[]byte{2, 0}, 512},
		{[]byte{4, 0}, 1024},
		{[]byte{255, 255, 255}, 0xFFFFFF},
	}
	for _, s := range cr {
		if r := toInt(s.in); !reflect.DeepEqual(s.out, r) {
			t.Errorf("got %v, expected %v\n", r, s.out)
		}
	}
}
