package ss

import (
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
		{"すもももももも", "pearch"},
		{"すもも", "pearch"},
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
		{"すもも", []string{"pearch", "もも"}},
		{"こんにちわ", nil},
		{"こんにちは", []string{"hello"}},
		{"世界", []string{"world"}},
		{"すもももももも", []string{"pearch"}},
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
