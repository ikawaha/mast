package main

import (
	"fmt"
	"github.com/ikawaha/mast/si"
)

func main() {
	pairs := si.PairSlice{
		{"こんにちは", 111},
		{"こんにちは", 222},
		{"こんばんは", 333},
	}

	t, _ := si.Build(pairs)
	vs := t.Search("こんにちは")
	fmt.Println(vs)
}
