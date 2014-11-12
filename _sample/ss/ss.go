package main

import (
	"fmt"
	"github.com/ikawaha/mast/ss"
)

func main() {
	pairs := ss.PairSlice{
		{"こんにちは", "hello"},
		{"こんにちは", "Здравствуйте"},
		{"こんばんは", "good evening"},
		{"東京", "Tokyo"},
		{"東京チョコレート", "Capsule"},
		{"東京チョコレート", "Eel"},
	}

	fst, _ := ss.Build(pairs)
	if o := fst.Search("こんにちは"); o != nil {
		fmt.Println(o)
	}

	inp := "東京チョコレートMIX"
	lens, outs := fst.CommonPrefixSearch(inp)
	for i := range outs {
		fmt.Println(inp[0:lens[i]], outs[i])
	}
}
