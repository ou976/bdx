package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/belldata-dx/bdx"
	"github.com/belldata-dx/bdx/interfaces"
)

type header struct {
	Key   string
	Value string
}

func request(r http.Handler, method, path, body string, headers ...header) *httptest.ResponseRecorder {
	headers = append(headers, header{
		Key:   "Content-Type",
		Value: "application/json",
	})
	reader := strings.NewReader(body)
	req := httptest.NewRequest(method, path, reader)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func setRouter() *bdx.Engine {
	router := bdx.New()
	router.GET("/", func(c interfaces.Context) {})
	return router
}
