package responder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func checkResponse(t *testing.T, resp *httptest.ResponseRecorder, code int, data string) {
	assert := assert.New(t)
	assert.Equal(resp.Body.String(), data)
	assert.Equal(resp.Code, code)
}

func checkJSON(t *testing.T, resp *httptest.ResponseRecorder) {
	assert := assert.New(t)
	assert.Equal(resp.Header().Get(elton.HeaderContentType), elton.MIMEApplicationJSON)
}

func checkContentType(t *testing.T, resp *httptest.ResponseRecorder, contentType string) {
	assert := assert.New(t)
	assert.Equal(resp.Header().Get(elton.HeaderContentType), contentType)
}

func TestResponder(t *testing.T) {
	m := New(Config{
		Fastest: true,
	})
	req := httptest.NewRequest("GET", "https://aslant.site/", nil)

	t.Run("skip", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		fn := New(Config{
			Skipper: func(c *elton.Context) bool {
				return true
			},
		})
		err := fn(c)
		assert.Nil(err)
		assert.True(done)
	})

	t.Run("return error", func(t *testing.T) {
		assert := assert.New(t)
		customErr := errors.New("abcd")
		c := elton.NewContext(nil, nil)
		done := false
		c.Next = func() error {
			done = true
			return customErr
		}
		fn := NewDefault()
		err := fn(c)
		assert.Equal(err, customErr)
		assert.True(done)
	})

	t.Run("set BodyBuffer", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		done := false
		c.Next = func() error {
			c.BodyBuffer = bytes.NewBuffer([]byte(""))
			done = true
			return nil
		}
		fn := New(Config{})
		err := fn(c)
		assert.Nil(err)
		assert.True(done)
	})

	t.Run("invalid response", func(t *testing.T) {
		d := elton.New()
		d.Use(m)
		d.GET("/", func(c *elton.Context) error {
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		checkResponse(t, resp, 500, "category=elton-responder, message=invalid response")
	})

	t.Run("return string", func(t *testing.T) {
		d := elton.New()
		d.Use(m)
		d.GET("/", func(c *elton.Context) error {
			c.Body = "abc"
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		checkResponse(t, resp, 200, "abc")
		checkContentType(t, resp, "text/plain; charset=UTF-8")
	})

	t.Run("return bytes", func(t *testing.T) {
		d := elton.New()
		d.Use(m)
		d.GET("/", func(c *elton.Context) error {
			c.Body = []byte("abc")
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		checkResponse(t, resp, 200, "abc")
		checkContentType(t, resp, elton.MIMEBinary)
	})

	t.Run("return struct", func(t *testing.T) {
		type T struct {
			Name string `json:"name,omitempty"`
		}
		d := elton.New()
		d.Use(m)
		d.GET("/", func(c *elton.Context) error {
			c.Created(&T{
				Name: "tree.xie",
			})
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		checkResponse(t, resp, 201, `{"name":"tree.xie"}`)
		checkJSON(t, resp)
	})

	t.Run("json marshal fail", func(t *testing.T) {
		assert := assert.New(t)
		d := elton.New()
		d.Use(m)
		d.GET("/", func(c *elton.Context) error {
			c.Body = func() {}
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		assert.Equal(500, resp.Code)
		assert.Equal("message=json: unsupported type: func()", resp.Body.String())
	})

	t.Run("reader body", func(t *testing.T) {
		assert := assert.New(t)
		d := elton.New()
		d.Use(m)
		d.GET("/", func(c *elton.Context) error {
			c.Body = bytes.NewReader([]byte("abcd"))
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		assert.Equal(resp.Code, 200)
		assert.Equal(resp.Body.String(), "abcd")
	})
}

type HelloWord struct {
	Content string  `json:"content,omitempty"`
	Size    int     `json:"size,omitempty"`
	Price   float32 `json:"price,omitempty"`
	VIP     bool    `json:"vip,omitempty"`
}

func getBenchmarkData() *HelloWord {
	arr := make([]string, 0)
	for i := 0; i < 100; i++ {
		arr = append(arr, "花褪残红青杏小。燕子飞时，绿水人家绕。枝上柳绵吹又少，天涯何处无芳草！")
	}
	content := strings.Join(arr, "\n")
	data := &HelloWord{
		Content: content,
		Size:    100,
		Price:   10.12,
		VIP:     true,
	}
	return data
}

func BenchmarkJSON(b *testing.B) {
	data := getBenchmarkData()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Fatalf("json marshal fail, %v", err)
		}
	}
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.9 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
