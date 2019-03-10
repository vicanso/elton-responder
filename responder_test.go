package responder

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vicanso/cod"
)

func checkResponse(t *testing.T, resp *httptest.ResponseRecorder, code int, data string) {
	if resp.Body.String() != data ||
		resp.Code != code {
		t.Fatalf("check response fail")
	}
}

func checkJSON(t *testing.T, resp *httptest.ResponseRecorder) {
	if resp.Header().Get(cod.HeaderContentType) != cod.MIMEApplicationJSON {
		t.Fatalf("response content type should be json")
	}
}

func checkContentType(t *testing.T, resp *httptest.ResponseRecorder, contentType string) {
	if resp.Header().Get(cod.HeaderContentType) != contentType {
		t.Fatalf("response content type check fail")
	}
}

func TestResponder(t *testing.T) {
	m := NewDefault()
	req := httptest.NewRequest("GET", "https://aslant.site/", nil)

	t.Run("skip", func(t *testing.T) {
		c := cod.NewContext(nil, nil)
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		fn := New(Config{
			Skipper: func(c *cod.Context) bool {
				return true
			},
		})
		err := fn(c)
		if err != nil ||
			!done {
			t.Fatalf("skip fail")
		}
	})

	t.Run("set BodyBuffer", func(t *testing.T) {
		c := cod.NewContext(nil, nil)
		done := false
		c.Next = func() error {
			c.BodyBuffer = bytes.NewBuffer([]byte(""))
			done = true
			return nil
		}
		fn := New(Config{})
		err := fn(c)
		if err != nil ||
			!done {
			t.Fatalf("set body buffer should pass")
		}
	})

	t.Run("invalid response", func(t *testing.T) {
		d := cod.New()
		d.Use(m)
		d.GET("/", func(c *cod.Context) error {
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		checkResponse(t, resp, 500, "category=cod-responder, message=invalid response")
	})

	t.Run("return string", func(t *testing.T) {
		d := cod.New()
		d.Use(m)
		d.GET("/", func(c *cod.Context) error {
			c.SetContentTypeByExt(".html")
			c.Body = "abc"
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		checkResponse(t, resp, 200, "abc")
		checkContentType(t, resp, "text/html; charset=utf-8")
	})

	t.Run("return bytes", func(t *testing.T) {
		d := cod.New()
		d.Use(m)
		d.GET("/", func(c *cod.Context) error {
			c.Body = []byte("abc")
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		checkResponse(t, resp, 200, "abc")
		checkContentType(t, resp, cod.MIMEBinary)
	})

	t.Run("return struct", func(t *testing.T) {
		type T struct {
			Name string `json:"name,omitempty"`
		}
		d := cod.New()
		d.Use(m)
		d.GET("/", func(c *cod.Context) error {
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
		d := cod.New()
		d.Use(m)
		d.GET("/", func(c *cod.Context) error {
			c.Body = func() {}
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		checkResponse(t, resp, 500, `{"statusCode":500,"message":"func() is unsupported type","exception":true}`)
	})
}

type HelloWord struct {
	Content string  `json:"content,omitempty"`
	Size    int     `json:"size,omitempty"`
	Price   float32 `json:"price,omitempty"`
	VIP     bool    `json:"vip,omitempty"`
}

func BenchmarkJSON(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Fatalf("json marshal fail, %v", err)
		}
	}
}

func BenchmarkStandardJSON(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		_, err := standardJSON.Marshal(data)
		if err != nil {
			b.Fatalf("standard json marshal fail, %v", err)
		}
	}
}

func BenchmarkFastJSON(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		_, err := fastJSON.Marshal(data)
		if err != nil {
			b.Fatalf("fast json marshal fail, %v", err)
		}
	}
}
