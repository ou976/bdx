package bdx

import (
	contextPac "context"
	"net/http"

	logger "github.com/belldata-dx/bdx-logger"
	"github.com/belldata-dx/bdx/bdxctx"
	"github.com/belldata-dx/bdx/interfaces"
	"github.com/belldata-dx/bdx/middleware"
	"github.com/belldata-dx/bdx/param"
	"github.com/julienschmidt/httprouter"
)

const (
	defaultMaxMultipartMemory = 32 << 20 // 32MB
)

type (
	contextKey struct{}

	// Engine bdxフレームワークインスタンス
	Engine struct {
		RouterGroup
		maxMultipartMemory int64
		maxParams          uint16
		log                logger.ILogger
	}
)

// DefaultLogger デフォルトで使用されるlogger
var DefaultLogger logger.ILogger

// ContextKey `request.Context()`で取得できるコンテキストキー
var ContextKey = contextKey{}

var _ interfaces.Engine = &Engine{}

// New bdx Instance
func New() (engine *Engine) {
	DefaultLogger = logger.New("engine", logger.Info)
	engine = &Engine{
		RouterGroup: RouterGroup{
			root:     true,
			route:    httprouter.New(),
			basePath: "/",
		},
		log:                DefaultLogger,
		maxMultipartMemory: defaultMaxMultipartMemory,
	}
	engine.engine = engine
	engine.pool.New = func() interface{} {
		v := make(param.Params, 0, engine.maxParams)
		return bdxctx.New(engine, &v)
	}
	return
}

// Default はデフォルト設定されたルータ
func Default() (engine *Engine) {
	engine = New()
	engine.Use(middleware.Logger)
	return
}

// Use ミドルウェア追加
func (engine *Engine) Use(handlers ...interfaces.BdxHandlerFunc) interfaces.Routes {
	engine.RouterGroup.Use(handlers...)
	return engine
}

// ServeHTTP .
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := engine.pool.Get().(*bdxctx.Context)
	c.SetLogger(engine.log)
	c.Reset(w, r)
	handler := func(c interfaces.Context) {
		r := c.Request()
		ctx := contextPac.WithValue(r.Context(), ContextKey, c)
		r = r.WithContext(ctx)
		engine.RouterGroup.route.ServeHTTP(c.Response(), r)
	}
	position := 0
	handlers := engine.RouterGroup.Handlers
	if len(handlers) > 0 {
		handlers = append(handlers[:position+1], handlers[position:]...)
		handlers[position] = handler
	} else {
		handlers = append(handlers, handler)
	}
	c.SetHandler(handlers)
	c.Next()
	engine.pool.Put(c)
}

// Run ListenAndServe
func (engine *Engine) Run(addr ...string) error {
	address := resolveAddress(addr)
	return http.ListenAndServe(address, engine)
}

// RunTLS ListenAndServeTLS
func (engine *Engine) RunTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, engine)
}

// SetLogger Logger Change
func (engine *Engine) SetLogger(log logger.ILogger) {
	DefaultLogger = log
	engine.log = DefaultLogger
}

// SetLogLevel ログ出力レベルを設定
func (engine *Engine) SetLogLevel(level logger.LogLevel) {
	engine.log.SetLevel(level)
}

// MaxMultipartMemory Multipartコンテンツタイプで許容される量
func (engine *Engine) MaxMultipartMemory() int64 {
	return engine.maxMultipartMemory
}
