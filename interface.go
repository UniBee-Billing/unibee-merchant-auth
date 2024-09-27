package unibee_merchant_auth

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IContext interface {
	Init(r *ghttp.Request, customCtx *UniBeeContext)
	GetUniBeeContext(ctx context.Context) *UniBeeContext
	SetUser(ctx context.Context, ctxUser *UniBeeContextUser)
	SetMerchantMember(ctx context.Context, ctxMerchantMember *UniBeeContextMerchantMember)
	SetData(ctx context.Context, data g.Map)
}

var singleTonContext IContext

func Context() IContext {
	if singleTonContext == nil {
		panic("implement not found for interface IContext, forgot register?")
	}
	return singleTonContext
}

const (
	SystemAssertPrefix = "system_assert: "
)

func GetMerchantId(ctx context.Context) uint64 {
	if Context().GetUniBeeContext(ctx) == nil {
		panic(SystemAssertPrefix + "Context Not Found")
	}
	if Context().GetUniBeeContext(ctx).MerchantId <= 0 {
		panic(SystemAssertPrefix + "Invalid Merchant")
	}
	return Context().GetUniBeeContext(ctx).MerchantId
}

func RegisterContext(i IContext) {
	singleTonContext = i
}
