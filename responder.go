// Copyright 2018 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package responder

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
)

type (
	// Config response config
	Config struct {
		Skipper elton.Skipper
		// Fastest set to true will use fast json
		Fastest bool
		// Marshal custom marshal function
		Marshal func(v interface{}) ([]byte, error)
		// ContentType response's content type
		ContentType string
	}
)

const (
	// ErrCategory responder error category
	ErrCategory = "elton-responder"
)

var (
	// errInvalidResponse invalid response(body an status is nil)
	errInvalidResponse = &hes.Error{
		Exception:  true,
		StatusCode: 500,
		Message:    "invalid response",
		Category:   ErrCategory,
	}
)

// NewDefault create a default responder
func NewDefault() elton.Handler {
	return New(Config{})
}

// New create a responder
func New(config Config) elton.Handler {
	skipper := config.Skipper
	if skipper == nil {
		skipper = elton.DefaultSkipper
	}
	marshal := config.Marshal
	// 如果未定义marshal
	if marshal == nil {
		marshal = json.Marshal
	}
	contentType := config.ContentType
	if contentType == "" {
		contentType = elton.MIMEApplicationJSON
	}

	return func(c *elton.Context) (err error) {
		if skipper(c) {
			return c.Next()
		}
		err = c.Next()
		if err != nil {
			return
		}
		// 如果已设置了BodyBuffer，则已生成好响应数据，跳过
		if c.BodyBuffer != nil {
			return
		}

		if c.StatusCode == 0 && c.Body == nil {
			// 如果status code 与 body 都为空，则为非法响应
			err = errInvalidResponse
			return
		}
		// 如果body是reader，则跳过
		if c.IsReaderBody() {
			return
		}

		ct := elton.HeaderContentType

		hadContentType := false
		// 判断是否已设置响应头的Content-Type
		if c.GetHeader(ct) != "" {
			hadContentType = true
		}

		var body []byte
		if c.Body != nil {
			switch data := c.Body.(type) {
			case string:
				if !hadContentType {
					c.SetHeader(ct, elton.MIMETextPlain)
				}
				body = []byte(data)
			case []byte:
				if !hadContentType {
					c.SetHeader(ct, elton.MIMEBinary)
				}
				body = data
			default:
				// 转换为json
				buf, e := marshal(data)
				if e != nil {
					he := hes.NewWithErrorStatusCode(e, http.StatusInternalServerError)
					he.Exception = true
					err = he
					return
				}
				if !hadContentType {
					c.SetHeader(ct, contentType)
				}
				body = buf
			}
		}

		statusCode := c.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		if len(body) != 0 {
			c.BodyBuffer = bytes.NewBuffer(body)
		}
		c.StatusCode = statusCode
		return nil
	}
}
