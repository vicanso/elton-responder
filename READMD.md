# cod-responder

[![Build Status](https://img.shields.io/travis/vicanso/cod-responder.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-responder)

Responder middleware for cod, it can convert `Context.Body` to json data. By this middleware, it's more simple for successful response.


```go
package main

import (
	"github.com/vicanso/cod"

	responder "github.com/vicanso/cod-responder"
)

func main() {
	d := cod.New()

	d.Use(responder.NewDefault())

	// {"name":"tree.xie","id":123}
	d.GET("/", func(c *cod.Context) (err error) {
		c.Body = &struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		}{
			"tree.xie",
			123,
		}
		return
	})

	d.ListenAndServe(":7001")
}
```