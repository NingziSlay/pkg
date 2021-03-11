# components

![Test](https://github.com/NingziSlay/components/workflows/Test/badge.svg) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/zerolog/master/LICENSE) [![Coverage](http://gocover.io/_badge/github.com/NingziSlay/components)](https://gocover.io/github.com/NingziSlay/components) [![Build Status](https://travis-ci.com/NingziSlay/components.svg?branch=main)](https://travis-ci.com/github/NingziSlay/components)
## install

```shell
go get -u github.com/NingziSlay/components
```

## config

usage:

```golang
package main

import (
	"github.com/NingziSlay/components"
	"log"
)

type Config struct{}

func main() {
	var c Config
	if err := components.MustMapConfig(&c); err != nil {
		log.Fatal(err)
	}
	// use c here
}
```