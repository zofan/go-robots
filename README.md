[![Build Status](https://travis-ci.org/zofan/go-robots.svg?branch=master)](https://travis-ci.org/zofan/go-robots)
[![Go Report Card](https://goreportcard.com/badge/github.com/zofan/go-robots)](https://goreportcard.com/report/github.com/zofan/go-robots)
[![Coverage Status](https://coveralls.io/repos/github/zofan/go-robots/badge.svg?branch=master)](https://coveralls.io/github/zofan/go-robots?branch=master)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/zofan/go-robots)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/zofan/go-robots/master/LICENSE)

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