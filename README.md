[![Go Report Card](https://goreportcard.com/badge/github.com/zofan/go-robots)](https://goreportcard.com/report/github.com/zofan/go-robots)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/zofan/go-robots)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/zofan/go-robots/master/LICENSE)
[![Sourcegraph](https://sourcegraph.com/github.com/zofan/go-robots/-/badge.svg)](https://sourcegraph.com/github.com/zofan/go-robots?badge)
[![Code Climate](https://codeclimate.com/github/zofan/go-robots/badges/gpa.svg)](https://codeclimate.com/github/zofan/go-robots)
[![Test Coverage](https://codeclimate.com/github/zofan/go-robots/badges/coverage.svg)](https://codeclimate.com/github/zofan/go-robots)
[![HitCount](http://hits.dwyl.io/zofan/go-robots.svg)](http://hits.dwyl.io/zofan/go-robots)

#### Features
- Easy and minimalistic

#### Install

> go get -u github.com/zofan/go-robots

#### Usage example

```$go
package main

import (
	"github.com/zofan/go-robots"
	"net/http"
)

func main() {
	client := &http.Client{}
	req, _ := http.NewRequest(`GET`, `https://www.bbc.com/robots.txt`, nil)

	resp, _ := client.Do(req)
	if resp == nil {
		return
	}

	config, err := robots.ParseResponse(resp)
	if err == nil {
	    // config processing
	}
}
```