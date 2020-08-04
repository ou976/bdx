package bdx

import (
	"net/http"

	"github.com/belldata-dx/bdx/interfaces"
	"github.com/belldata-dx/bdx/param"
	"github.com/julienschmidt/httprouter"
)

// Handle はハンドラの型
type Handle func(http.ResponseWriter, *http.Request, param.IParam)

// HandleFunc は`Handle`を`BdxHandlerFunc`へ変換
func HandleFunc(h Handle) interfaces.BdxHandlerFunc {
	return func(c interfaces.Context) {
		h(c.Response(), c.Request(), c.Params())
	}
}

// HandlerFunc は`http.HandlerFunc`を`BdxHandlerFunc`へ変換
func HandlerFunc(h http.HandlerFunc) interfaces.BdxHandlerFunc {
	return func(c interfaces.Context) {
		h.ServeHTTP(c.Response(), c.Request())
	}
}

// HandleFuncrouter は`httprouter.Handle`を`BdxHandlerFunc`へ変換
func HandleFuncrouter(h httprouter.Handle) interfaces.BdxHandlerFunc {
	return func(c interfaces.Context) {
		cParams := c.Params()
		params := httprouter.Params{}
		for _, cParam := range cParams {
			param := httprouter.Param{
				Key:   cParam.Key,
				Value: cParam.Value,
			}
			params = append(params, param)
		}
		h(c.Response(), c.Request(), params)
	}
}
