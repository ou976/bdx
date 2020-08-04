package bdx_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/belldata-dx/bdx"
	logger "github.com/belldata-dx/bdx-logger"
	"github.com/belldata-dx/bdx/interfaces"
	"github.com/belldata-dx/bdx/middleware"
	"github.com/belldata-dx/bdx/param"
)

var (
	log = logger.New("bdx_test", logger.Info)
)

type header struct {
	Key   string
	Value string
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Debug(err)
	}
	json.Unmarshal(body, &u)
	responseBody, _ := json.Marshal(&u)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(responseBody))
}

func index(w http.ResponseWriter, r *http.Request, pm param.IParam) {
	fmt.Fprint(w, "aaa")
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

func newlog() logger.ILogger {
	return logger.New("test", logger.Debug)
}

func TestRouter(t *testing.T) {
	var w *httptest.ResponseRecorder
	var req *http.Request
	var read []byte
	router := bdx.New()
	router.SetLogger(newlog())
	router.Use(middleware.Logger)
	router.GET("/", bdx.HandleFunc(index))

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	read, _ = ioutil.ReadAll(w.Body)
	assert.Equal(t, "aaa", string(read))
}

func TestGroup(t *testing.T) {
	var w *httptest.ResponseRecorder
	var read []byte

	router := bdx.Default()
	router.SetLogger(newlog())
	g := router.Group("/v1")
	{
		g.Use(middleware.Validator(User{}))
		g.POST("/user", bdx.HandlerFunc(userHandler))
		g.POST("/users", bdx.HandlerFunc(userHandler))
	}

	jsonBody := `{"id": 1, "name": ""}`
	w = request(router, http.MethodPost, "/v1/user", jsonBody)
	read, _ = ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"code":400,"error":"Invalid body parser","error_descript":"Invalid body parser","error_detail":{"Name":"Nameは必須フィールドです"}}`, string(read))

	jsonBody = `{"id":1,"name":"aaaa"}`
	w = request(router, http.MethodPost, "/v1/users", jsonBody)
	read, _ = ioutil.ReadAll(w.Body)
	assert.Equal(t, jsonBody, string(read))

}

func TestPathParameter(t *testing.T) {
	var id string
	router := bdx.Default()
	router.SetLogger(newlog())
	router.GET("/:id", func(c interfaces.Context) {
		p := c.Params()
		id = p.ByName("id")
	})
	request(router, http.MethodGet, "/123", "")
	assert.Equal(t, id, "123")
}
