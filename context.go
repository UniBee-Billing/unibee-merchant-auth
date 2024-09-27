package unibee_merchant_auth

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
)

type SContext struct{}

func init() {
	RegisterContext(New())
}

func New() *SContext {
	return &SContext{}
}

func (s *SContext) Init(r *ghttp.Request, uniBeeContext *UniBeeContext) {
	r.SetCtxVar(ContextKey, uniBeeContext)
}

func (s *SContext) Get(ctx context.Context) *UniBeeContext {
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*UniBeeContext); ok {
		return localCtx
	}
	return nil
}
