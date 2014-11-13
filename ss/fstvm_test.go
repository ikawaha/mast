package ss

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestFstVMSearch01(t *testing.T) {
	inp := PairSlice{
		{"1a22xss", "world"},
		{"1b22yss", "goodby"},
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

func TestSaveAndLoad01(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", "hello"},
		{"世界", "world"},
		{"すもももももも", "peach"},
		{"すもも", "peach"},
		{"すもも", "もも"},
	}
	v1, e := Build(inp)
	if e != nil {
		t.Errorf("unexpected error %v\n", e)
	}
	var b bytes.Buffer
	if e := v1.Save(&b); e != nil {
		t.Errorf("unexpected error %v\n", e)
	}
	v2 := FstVM{}
	v2.Load(&b)
	if !reflect.DeepEqual(v1, v2) {
		t.Errorf("save:\n%v\nload:\n%v\n", v1, v2)
	}

}
