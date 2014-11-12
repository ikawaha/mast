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
		{"東京", 444},
		{"東京チョコレート", 555},
		{"東京チョコレート", 666},
	}

	fst, _ := si.Build(pairs)
	if o := fst.Search("こんにちは"); o != nil {
		fmt.Println(o)
	}
	inp := "東京チョコレートMIX"
	lens, outs := fst.CommonPrefixSearch(inp)
	for i := range outs {
		fmt.Println(inp[0:lens[i]], outs[i])
	}

}
