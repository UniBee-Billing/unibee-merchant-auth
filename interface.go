package unibee_merchant_auth

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IContext interface {
	Init(r *ghttp.Request, uniBeeContext *UniBeeContext)
	Get(ctx context.Context) *UniBeeContext
}

var singleTonContext IContext

func Context() IContext {
	if singleTonContext == nil {
		panic("implement not found for interface IContext, forgot register?")
	}
	return singleTonContext
}

func RegisterContext(i IContext) {
	singleTonContext = i
}

const (
	SystemAssertPrefix = "system_assert: "
)
