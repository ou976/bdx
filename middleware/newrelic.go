package middleware

import (
	"context"

	"github.com/belldata-dx/bdx/interfaces"
	newrelic "github.com/newrelic/go-agent"
)

type newrelicKey struct{}

var (
	// NewrelicKey Context key
	NewrelicKey = newrelicKey{}
)

// NewrelicMiddleware newrelic monitoringツールを使用する際のミドルウェア
func NewrelicMiddleware(app newrelic.Application) interfaces.BdxHandlerFunc {
	return func(ctx interfaces.Context) {
		req := ctx.Request()
		res := ctx.Response()
		tx := app.StartTransaction(req.URL.Path, res, req)
		defer tx.End()
		cont := context.WithValue(req.Context(), NewrelicKey, tx)
		ctx.SetRequest(req.WithContext(cont))
		ctx.Next()
	}
}
