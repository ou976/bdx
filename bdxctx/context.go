package bdxctx

import (
	"math"
	"net/http"
	"net/url"

	logger "github.com/belldata-dx/bdx-logger"
	"github.com/belldata-dx/bdx/interfaces"
	"github.com/belldata-dx/bdx/param"
	"github.com/belldata-dx/bdx/render"
)

const abortIndex int8 = math.MaxInt8 / 2

type (
	// Context はハンドラへ渡される属性
	Context struct {
		request    *http.Request
		response   http.ResponseWriter
		handlers   interfaces.HandlersChain
		index      int8
		err        error
		engine     interfaces.Engine
		params     *param.Params
		logger     logger.ILogger
		queryCache url.Values
		formCache  url.Values
	}
)

var _ interfaces.Context = &Context{}

// New Context Constructor
func New(engine interfaces.Engine, params *param.Params) *Context {
	return &Context{engine: engine, params: params}
}

// Request `*http.Request`
func (c *Context) Request() *http.Request {
	return c.request
}

// SetRequest `*http.Request`
func (c *Context) SetRequest(req *http.Request) {
	c.request = req
}

// Response `http.ResponseWriter`
func (c *Context) Response() http.ResponseWriter {
	return c.response
}

// Logger `logger.Logger`
func (c *Context) Logger() logger.ILogger {
	return c.logger
}

// SetLogger `logger.Logger`
func (c *Context) SetLogger(log logger.ILogger) {
	c.logger = log
}

// Reset .
func (c *Context) Reset(response http.ResponseWriter, request *http.Request) {
	c.response = response
	c.request = request
	c.index = -1
	c.handlers = interfaces.HandlersChain{}
	*c.params = (*c.params)[0:0]
	c.err = nil
}

// Next は次のミドルウェアもしくはハンドラを実行
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// IsAborted はどこかで終了処理が呼び出されたかどうか
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort 後続を処理せず`この処理`で終了
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatusAndMessage 後続を処理せず`この処理`で終了し、エラーレスポンスを生成
func (c *Context) AbortWithStatusAndMessage(status int, buf []byte) {
	w := c.response
	w.WriteHeader(status)
	w.Write(buf)
	c.Abort()
}

// AbortWithUnsupportedMediaType 後続を処理せず`この処理`で終了し、エラーレスポンスを生成
func (c *Context) AbortWithUnsupportedMediaType() {
	status := http.StatusUnsupportedMediaType
	w := c.response
	w.WriteHeader(status)
	c.Abort()
}

// Params URIパス パラメータ
func (c *Context) Params() param.Params {
	return *c.params
}

// SetParams paramsをセット
func (c *Context) SetParams(params *param.Params) {
	c.params = params
}

// SetHandler ハンドラー設定
func (c *Context) SetHandler(handlers interfaces.HandlersChain) {
	c.handlers = handlers
}

// Handler 実際に登録されているハンドラ
func (c *Context) Handler() interfaces.BdxHandlerFunc {
	return c.handlers.Last()
}

// Status HTTP response codeを設定
func (c *Context) Status(code int) {
	c.response.WriteHeader(code)
}

// Render HTTTP response codeと`render.Render`にrender dataを書き込みます
func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.response)
		return
	}

	if err := r.Render(c.response); err != nil {
		panic(err)
	}
}

// JSON JSONでHTTP responseを書き込み
func (c *Context) JSON(code int, data interface{}) {
	c.Render(code, render.JSON{Data: data})
}

// XML XMLでHTTP responseを書き込み
func (c *Context) XML(code int, data interface{}) {
	c.Render(code, render.XML{Data: data})
}

// YAML YAMLでHTTP responseを書き込み
func (c *Context) YAML(code int, data interface{}) {
	c.Render(code, render.YAML{Data: data})
}

// Query url parameterが存在すればそれを返します。
// 存在しない場合は空文字を返します。
// これは`c.Request.URL.Query().Get(key)`のショートカットと同じです。
//     GET /path?id=1&name=hoge&value=
//     c.Query("id") == "1"
//     c.Query("name") == "hogehoe"
//     c.Query("value") == ""
//     c.Query("foo") == ""
func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

// DefaultQuery url parameterが存在すればそれを返します。
// 存在しない場合は`defaultValue`に設定された値を返します。
//     GET /path?id=1&name=hoge&value=
//     c.DefaultQuery("id", "123") == "1"
//     c.DefaultQuery("name", "none") == "hoge"
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

// GetQueryArray url parameterが存在すれば複数の値をを返し、さらに`true`も返します。
// 存在しない場合は空の文字列配列と`false`を返します。
func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQuery()
	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// QueryArray url parameterが存在すれば複数の値をを返します。
// 存在しない場合は空の文字列配列を返します。
func (c *Context) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

// GetQuery url parameterが存在すればそれを返し、さらに`true`を返します。
// 存在しない場合は空文字を返し、`false`も返します。
func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], true
	}

	return "", false
}

// initQuery `queryCache`へ値を移送する処理。
// `c.Reqest.URL.Query()`
func (c *Context) initQuery() {
	if c.queryCache == nil {
		c.queryCache = c.request.URL.Query()
	}
}

// initForm `formCache`へ値を移送する処理。
func (c *Context) initPostForm() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		req := c.request
		if err := req.ParseMultipartForm(c.engine.MaxMultipartMemory()); err != nil {
			c.logger.Debugf("マルチパート形式の配列を解析する際のエラー: %v", err)
		}
		c.formCache = c.request.PostForm
	}
}

// PostForm POST urlencoded form or multipart formが存在すればそれを返します。
// 存在しない場合は空文字を返します。
//     POST /path
//     Content-Type: application/x-www-from-urlencoded
//
//     id=1&name=hoge
//     c.PostForm("id") == "1"
//     c.PostForm("name") == "hogehoe"
//     c.PostForm("value") == ""
//     c.PostForm("foo") == ""
func (c *Context) PostForm(key string) string {
	value, _ := c.GetPostForm(key)
	return value
}

// DefaultPostForm POST urlencoded form or multipart formが存在すればそれを返します。
// 存在しない場合は`defaultValue`に設定された値を返します。
//     POST /path
//     Content-Type: application/x-www-from-urlencoded
//
//     name=hoge
//     c.DefaultPostForm("id", "123") == "123"
//     c.DefaultPostForm("name", "none") == "hoge"
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value, ok := c.GetPostForm(key); ok {
		return value
	}
	return defaultValue
}

// GetPostFormArray POST urlencoded form or multipart formが存在すれば複数の値をを返し、さらに`true`も返します。
// 存在しない場合は空の文字列配列と`false`を返します。
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.initPostForm()
	if values, ok := c.formCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// PostFormArray POST urlencoded form or multipart formが存在すれば複数の値をを返します。
// 存在しない場合は空の文字列配列を返します。
func (c *Context) PostFormArray(key string) []string {
	values, _ := c.GetPostFormArray(key)
	return values
}

// GetPostForm POST urlencoded form or multipart formすればそれを返し、さらに`true`を返します。
// 存在しない場合は空文字を返し、`false`も返します。
func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], true
	}

	return "", false
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
