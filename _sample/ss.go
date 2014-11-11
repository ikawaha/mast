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
	}

	t, _ := ss.Build(pairs)
	gs := t.Search("こんにちは")
	for _, g := range gs {
		fmt.Println(g)
	}
}
