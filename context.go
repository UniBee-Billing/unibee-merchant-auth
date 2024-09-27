package unibee_merchant_auth

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type SContext struct{}

func init() {
	RegisterContext(New())
}

func New() *SContext {
	return &SContext{}
}

func (s *SContext) Init(r *ghttp.Request, customCtx *UniBeeContext) {
	r.SetCtxVar(ContextKey, customCtx)
}

func (s *SContext) GetUniBeeContext(ctx context.Context) *UniBeeContext {
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*UniBeeContext); ok {
		return localCtx
	}
	return nil
}

func (s *SContext) SetUser(ctx context.Context, ctxUser *UniBeeContextUser) {
	s.GetUniBeeContext(ctx).User = ctxUser
}

func (s *SContext) SetMerchantMember(ctx context.Context, ctxMerchantMember *UniBeeContextMerchantMember) {
	s.GetUniBeeContext(ctx).MerchantMember = ctxMerchantMember
}

func (s *SContext) SetData(ctx context.Context, data g.Map) {
	s.GetUniBeeContext(ctx).Data = data
}
