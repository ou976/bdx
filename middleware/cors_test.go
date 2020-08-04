package middleware_test

import (
	"net/http"
	"testing"

	"github.com/belldata-dx/bdx/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	router := setRouter()
	router.Use(middleware.CORS)
	w := request(router, http.MethodGet, "/", "")
	acao := w.Header().Get("Access-Control-Allow-Origin")
	acah := w.Header().Get("Access-Control-Allow-Headers")
	assert.Equal(t, acao, "*")
	assert.Equal(t, acah, "*")
}
