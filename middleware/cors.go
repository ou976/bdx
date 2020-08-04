package middleware

import (
	"github.com/belldata-dx/bdx/interfaces"
)

// CORS CORSの適応
func CORS(ctx interfaces.Context) {
	w := ctx.Response()
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	ctx.Next()
}
