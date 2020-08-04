package middleware

import (
	"net/http"
	"reflect"
	"unsafe"

	"github.com/belldata-dx/bdx/interfaces"
)

// Logger ログ出力
func Logger(ctx interfaces.Context) {
	r := ctx.Request()
	w := ctx.Response()
	bdxlog := ctx.Logger()
	rAddr := r.RemoteAddr
	method := r.Method
	path := r.URL.Path
	bdxlog.Infof("Remote: %s [%s] %s", rAddr, method, path)
	ctx.Next()
	if _, ok := w.(http.ResponseWriter); ok {
		rv := reflect.ValueOf(w).Elem()
		piv := rv.FieldByName("status")
		val := reflect.Value{}
		if piv != val {
			pi := (*int)(unsafe.Pointer(piv.UnsafeAddr()))
			bdxlog.Infof("Status: %v", *pi)
		} else {
			bdxlog.Debug("http test")
		}
	}
	ctx.Next()
}
