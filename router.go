package bdx

import (
	"math"
	"net/http"
	"sync"

	"github.com/belldata-dx/bdx/bdxctx"
	"github.com/belldata-dx/bdx/interfaces"
	"github.com/belldata-dx/bdx/param"
	"github.com/julienschmidt/httprouter"
)

type (

	// RouterGroup ルータ
	RouterGroup struct {
		Handlers    interfaces.HandlersChain
		middlewares interfaces.HandlersChain
		pool        sync.Pool
		route       *httprouter.Router
		engine      *Engine
		root        bool
		basePath    string
	}
)

const abortIndex int8 = math.MaxInt8 / 2

var _ interfaces.Router = &RouterGroup{}

// Group ルータグループ
// `New()`した際は`/`がbasePath
func (group *RouterGroup) Group(relativePath string, handlers ...interfaces.BdxHandlerFunc) interfaces.Router {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
		route:    group.route,
	}
}

// Handler 新しいリクエストハンドルとミドルウェアを与えられたパスとメソッドで登録します。
// 最後のハンドルが実際のハンドルとして登録されます。
// それ以外はミドルウェアでなければいけません。
func (group *RouterGroup) Handler(method, relativePath string, handlers ...interfaces.BdxHandlerFunc) interfaces.Routes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.route.Handle(method, absolutePath, func(w http.ResponseWriter, rq *http.Request, pm httprouter.Params) {
		c := rq.Context().Value(ContextKey).(*bdxctx.Context)
		c.Reset(w, rq)
		params := param.Params{}
		for _, p := range pm {
			param := param.Param{
				Key:   p.Key,
				Value: p.Value,
			}
			params = append(params, param)
		}
		c.SetParams(&params)
		// c.params = &params
		c.Params()
		c.SetHandler(append(group.middlewares, handlers...))
		c.Next()
	})
	return group.returnObj()
}

// GET は`router.Handler("GET", path, handle)`のショートカットです。
func (group *RouterGroup) GET(relativePath string, handlers ...interfaces.BdxHandlerFunc) interfaces.Routes {
	group.Handler(http.MethodGet, relativePath, handlers...)
	return group.returnObj()
}

// POST は`router.Handler("POST", path, handle)`のショートカットです。
func (group *RouterGroup) POST(relativePath string, handlers ...interfaces.BdxHandlerFunc) interfaces.Routes {
	group.Handler(http.MethodPost, relativePath, handlers...)
	return group.returnObj()
}

// PUT は`router.Handler("PUT", path, handle)`のショートカットです。
func (group *RouterGroup) PUT(relativePath string, handlers ...interfaces.BdxHandlerFunc) interfaces.Routes {
	group.Handler(http.MethodPut, relativePath, handlers...)
	return group.returnObj()
}

// DELETE は`router.Handler("DELETE", path, handle)`のショートカットです。
func (group *RouterGroup) DELETE(relativePath string, handlers ...interfaces.BdxHandlerFunc) interfaces.Routes {
	group.Handler(http.MethodDelete, relativePath, handlers...)
	return group.returnObj()
}

// OPTIONS は`router.Handler("OPTIONS", path, handle)`のショートカットです。
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...interfaces.BdxHandlerFunc) interfaces.Routes {
	group.Handler(http.MethodOptions, relativePath, handlers...)
	return group.returnObj()
}

// Any は`GET`,`POST`,`PUT`,`DELETE`のショートカットです。
func (group *RouterGroup) Any(relativePath string, handlers ...interfaces.BdxHandlerFunc) interfaces.Routes {
	group.GET(relativePath, handlers...)
	group.POST(relativePath, handlers...)
	group.PUT(relativePath, handlers...)
	group.DELETE(relativePath, handlers...)
	group.OPTIONS(relativePath, handlers...)
	return group.returnObj()
}

// Use ルータの後に実行されるミドルウェアを追加します。
// この機能はハンドラ全体に影響します。
func (group *RouterGroup) Use(middlewares ...interfaces.BdxHandlerFunc) interfaces.Routes {
	group.middlewares = append(group.middlewares, middlewares...)
	return group.returnObj()
}

func (group *RouterGroup) mergeMiddlewareHandler() interfaces.Routes {
	position := len(group.middlewares) - 1
	handlers := group.Handlers
	if len(handlers) > 0 {
		handlers = append(handlers[:position+1], handlers[position:]...)
		for i, middleware := range group.middlewares {
			handlers[i] = middleware
		}
	} else {
		for _, middleware := range group.middlewares {
			group.Handlers = append(group.Handlers, middleware)
		}
	}
	return group
}

func (group *RouterGroup) combineHandlers(handlers interfaces.HandlersChain) interfaces.HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(interfaces.HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func (group *RouterGroup) returnObj() interfaces.Routes {
	if group.root {
		return group.engine
	}
	return group
}
