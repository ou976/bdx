package interfaces

import (
	"net/http"

	logger "github.com/belldata-dx/bdx-logger"
	"github.com/belldata-dx/bdx/param"
	"github.com/belldata-dx/bdx/render"
)

type (

	// Router ルーターが提供する全機能
	Router interface {
		Routes
		Group(relativePath string, handlers ...BdxHandlerFunc) Router
	}

	// Routes はルータが提供する機能
	Routes interface {
		Handler(method string, path string, handlers ...BdxHandlerFunc) Routes
		GET(path string, handlers ...BdxHandlerFunc) Routes
		POST(path string, handlers ...BdxHandlerFunc) Routes
		PUT(path string, handlers ...BdxHandlerFunc) Routes
		DELETE(path string, handlers ...BdxHandlerFunc) Routes
		OPTIONS(path string, handlers ...BdxHandlerFunc) Routes
		Use(middleware ...BdxHandlerFunc) Routes
	}

	// Context はハンドラへ渡される属性
	Context interface {
		// Request `*http.Request`
		Request() *http.Request
		// SetRequest `*http.Request`
		SetRequest(*http.Request)
		// Response `http.ResponseWriter`
		Response() http.ResponseWriter
		// Logger `bdx-logger.Logger`
		Logger() logger.ILogger
		// Reset .
		Reset(response http.ResponseWriter, request *http.Request)
		// Next は次のミドルウェアもしくはハンドラを実行
		Next()
		// IsAborted はどこかで終了処理が呼び出されたかどうか
		IsAborted() bool
		// Abort 後続を処理せず`この処理`で終了
		Abort()
		// AbortWithStatusAndMessage 後続を処理せず`この処理`で終了し、エラーレスポンスを生成
		AbortWithStatusAndMessage(status int, buf []byte)
		// AbortWithUnsupportedMediaType 後続を処理せず`この処理`で終了し、エラーレスポンスを生成
		AbortWithUnsupportedMediaType()
		// Params URIパス パラメータ
		Params() param.Params
		// SetParams paramsをセット
		SetParams(params *param.Params)
		// SetHandler ハンドラー設定
		SetHandler(handlers HandlersChain)
		// Handler 実際に登録されているハンドラ
		Handler() BdxHandlerFunc
		// Status HTTP response codeを設定
		Status(code int)
		// Render HTTTP response codeと`render.Render`にrender dataを書き込みます
		Render(code int, r render.Render)
		// JSON JSONでHTTP responseを書き込み
		JSON(code int, data interface{})
		// XML XMLでHTTP responseを書き込み
		XML(code int, data interface{})
		// YAML YAMLでHTTP responseを書き込み
		YAML(code int, data interface{})
		// Query url parameterが存在すればそれを返します。
		// 存在しない場合は空文字を返します。
		// これは`c.Request.URL.Query().Get(key)`のショートカットと同じです。
		//     GET /path?id=1&name=hoge&value=
		//     c.Query("id") == "1"
		//     c.Query("name") == "hogehoe"
		//     c.Query("value") == ""
		//     c.Query("foo") == ""
		Query(key string) string
		// DefaultQuery url parameterが存在すればそれを返します。
		// 存在しない場合は`defaultValue`に設定された値を返します。
		//     GET /path?id=1&name=hoge&value=
		//     c.DefaultQuery("id", "123") == "1"
		DefaultQuery(key, defaultValue string) string
		// GetQueryArray url parameterが存在すれば複数の値をを返し、さらに`true`も返します。
		// 存在しない場合は空の文字列配列と`false`を返します。
		GetQueryArray(key string) ([]string, bool)
		// QueryArray url parameterが存在すれば複数の値をを返します。
		// 存在しない場合は空の文字列配列を返します。
		QueryArray(key string) []string
		// GetQuery url parameterが存在すればそれを返し、さらに`true`を返します。
		// 存在しない場合は空文字を返し、`false`も返します。
		GetQuery(key string) (string, bool)
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
		PostForm(key string) string
		// DefaultPostForm POST urlencoded form or multipart formが存在すればそれを返します。
		// 存在しない場合は`defaultValue`に設定された値を返します。
		//     POST /path
		//     Content-Type: application/x-www-from-urlencoded
		//
		//     name=hoge
		//     c.DefaultPostForm("id", "123") == "123"
		//     c.DefaultPostForm("name", "none") == "hoge"
		DefaultPostForm(key, defaultValue string) string
		// GetPostFormArray POST urlencoded form or multipart formが存在すれば複数の値をを返し、さらに`true`も返します。
		// 存在しない場合は空の文字列配列と`false`を返します。
		GetPostFormArray(key string) ([]string, bool)
		// PostFormArray POST urlencoded form or multipart formが存在すれば複数の値をを返します。
		// 存在しない場合は空の文字列配列を返します。
		PostFormArray(key string) []string
		// GetPostForm POST urlencoded form or multipart formすればそれを返し、さらに`true`を返します。
		// 存在しない場合は空文字を返し、`false`も返します。
		GetPostForm(key string) (string, bool)
	}

	// Engine bdxが提供する機能
	Engine interface {
		Routes
		MaxMultipartMemory() int64
	}

	// BdxHandlerFunc ハンドラ
	BdxHandlerFunc func(Context)

	// HandlersChain ハンドラチェーン
	HandlersChain []BdxHandlerFunc
)

// Last returns the last handler in the chain. ie. the last handler is the main one.
func (c HandlersChain) Last() BdxHandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
