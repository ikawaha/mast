# Minimal Acyclic Subsequential Transducers

[mast](http://github.com/ikawaha/mast) is a library for building of a finite state transducer called [minimal acyclic subsequential transducer](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=CD58961193540FBC807D500663EFD451?doi=10.1.1.24.3698&rep=rep1&type=pdf).

## Installation

```
go get github.com/ikawaha/mast/...
```

## Usage

### String to String Transducers

```
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
```
outputs
```
hello
Здравствуйте
```

### String to Integer Transducers

```
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
```

outputs
```
[111 222]
```

## References
* [Direct construction of minimal acyclic subsequential transducers](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=CD58961193540FBC807D500663EFD451?doi=10.1.1.24.3698&rep=rep1&type=pdf), Stoyan Mihov and Denis Maurel, 2001.
