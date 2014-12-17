[![Build Status](https://travis-ci.org/ikawaha/mast.svg?branch=master)](https://travis-ci.org/ikawaha/mast) [![Coverage Status](https://coveralls.io/repos/ikawaha/mast/badge.png?branch=master)](https://coveralls.io/r/ikawaha/mast?branch=master) [![experimental](http://badges.github.io/stability-badges/dist/experimental.svg)](http://github.com/badges/stability-badges)

# Minimal Acyclic Subsequential Transducers

[mast](http://github.com/ikawaha/mast) is a library for building of a finite state transducer called [minimal acyclic subsequential transducer](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=CD58961193540FBC807D500663EFD451?doi=10.1.1.24.3698&rep=rep1&type=pdf).

ex. 

|input|output|
|:---|:---|
|apr|30|
|aug|31|
|dec|31|
|feb|28, 29|

![Finite State Transducer](https://github.com/ikawaha/mast/wiki/images/fst_sample_01.png)

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
```

outputs

```
[hello Здравствуйте]
東京 [Tokyo]
東京チョコレート [Capsule Eel]
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
```

outputs

```
[222 111]
東京 [444]
東京チョコレート [555 666]
```

## References
* [Direct construction of minimal acyclic subsequential transducers](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=CD58961193540FBC807D500663EFD451?doi=10.1.1.24.3698&rep=rep1&type=pdf), Stoyan Mihov and Denis Maurel, 2001.
