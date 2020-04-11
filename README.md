# elton-responder

The middleware has been archived, please use the middleware of [elton](https://github.com/vicanso/elton).

[![Build Status](https://img.shields.io/travis/vicanso/elton-responder.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-responder)

Responder middleware for elton, it can convert `Context.Body` to json data. Using this middleware, it's more simple for successful response. More response type can be supported through custom marshal function and content type.


```go
package main

import (
	"github.com/vicanso/elton"

	responder "github.com/vicanso/elton-responder"
)

func main() {
	e := elton.New()

	e.Use(responder.NewDefault())

	// {"name":"tree.xie","id":123}
	e.GET("/", func(c *elton.Context) (err error) {
		c.Body = &struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		}{
			"tree.xie",
			123,
		}
		return
	})

	er := e.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	}
}
```

## API

- `Config.Skipper` skipper function to skip middleware
- `Config.Fastest` if set true will use the `json-iterator` fastest config for better performance, deprecated
- `Config.Marshal` custom marshal function
- `Config.ContentType` the coontent type for response
