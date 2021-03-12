# pkg

![Test](https://github.com/NingziSlay/pkg/workflows/Test/badge.svg)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/NingziSlay/pkg/master/LICENSE)
[![Coverage](http://gocover.io/_badge/github.com/NingziSlay/pkg)](https://gocover.io/github.com/NingziSlay/pkg)
[![Build Status](https://travis-ci.com/NingziSlay/pkg.svg?branch=main)](https://travis-ci.com/github/NingziSlay/pkg)
[![Go Report Card](https://goreportcard.com/badge/github.com/NingziSlay/pkg)](https://goreportcard.com/report/github.com/NingziSlay/pkg)

## install

```shell
go get -u github.com/NingziSlay/pkg
```

## config

TODO: support time.Time

usage:

```golang
package main

import (
	"github.com/NingziSlay/pkg"
	"log"
)

type Config struct{}

func main() {
	var c Config
	if err := pkg.MustMapConfig(&c); err != nil {
		log.Fatal(err)
	}
	// use c here
}
```