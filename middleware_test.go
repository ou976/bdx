package bdx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/belldata-dx/bdx/interfaces"
	"github.com/stretchr/testify/assert"
)

type header struct {
	Key   string
	Value string
}

func request(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestMiddleware(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c interfaces.Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.Use(func(c interfaces.Context) {
		signature += "C"
		c.Next()
	})
	router.GET("/", func(c interfaces.Context) {
		signature += "D"
	})
	request(router, "GET", "/")
	assert.Equal(t, "ACDB", signature)
}
