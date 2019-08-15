# elton-responder

[![Build Status](https://img.shields.io/travis/vicanso/elton-responder.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-responder)

Responder middleware for elton, it can convert `Context.Body` to json data. By this middleware, it's more simple for successful response.


```go
package main

import (
	"github.com/vicanso/elton"

	responder "github.com/vicanso/elton-responder"
)

func main() {
	d := elton.New()

	d.Use(responder.NewDefault())

	// {"name":"tree.xie","id":123}
	d.GET("/", func(c *elton.Context) (err error) {
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

## API

- `Config.Skipper` skipper function to skip middleware
- `Config.Fastest` if set true will use the `json-iterator` fastest config for better performance