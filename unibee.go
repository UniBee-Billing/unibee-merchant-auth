package unibee_merchant_auth

import (
	"context"
	"github.com/UniBee-Billing/unibee-merchant-auth/bean"
	"github.com/UniBee-Billing/unibee-merchant-auth/jwt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type OpenApiConfig struct {
	Id                      uint64 `json:"id"                      description:""`                         //
	Qps                     int    `json:"qps"                     description:"total qps control"`        // total qps control
	MerchantId              uint64 `json:"merchantId"              description:"merchant id"`              // merchant id
	Hmac                    string `json:"hmac"                    description:"webhook hmac key"`         // webhook hmac key
	Callback                string `json:"callback"                description:"callback url"`             // callback url
	ApiKey                  string `json:"apiKey"                  description:"api key"`                  // api key
	Token                   string `json:"token"                   description:"api token"`                // api token
	IsDeleted               int    `json:"isDeleted"               description:"0-UnDeleted，1-Deleted"`    // 0-UnDeleted，1-Deleted
	ValidIps                string `json:"validIps"                description:""`                         //
	GatewayCallbackResponse string `json:"gatewayCallbackResponse" description:"callback return response"` // callback return response
	CompanyId               int64  `json:"companyId"               description:"company id"`               // company id
}

type UniBeeContext struct {
	Session        *ghttp.Session
	MerchantId     uint64
	User           *UniBeeContextUser
	MerchantMember *bean.MerchantMember
	Merchant       *bean.Merchant
	RequestId      string
	Data           g.Map
	OpenApiConfig  *OpenApiConfig
	OpenApiKey     string
	IsOpenApiCall  bool
	Language       string
	UserAgent      string
	Authorization  string
	TokenString    string
	Token          *jwt.TokenClaims
}

type UniBeeContextUser struct {
	Id         uint64
	MerchantId uint64
	Token      string
	Email      string
	Lang       string
}

type UniBeeContextMerchantMember struct {
	Id         uint64
	MerchantId uint64
	Token      string
	Email      string
	IsOwner    bool
}

func GetUniBeeMerchantId(ctx context.Context) uint64 {
	if Context().Get(ctx) == nil {
		panic(SystemAssertPrefix + "Context Not Found")
	}
	if Context().Get(ctx).MerchantId <= 0 {
		panic(SystemAssertPrefix + "Invalid Merchant")
	}
	return Context().Get(ctx).MerchantId
}

func GetUniBeeContext(ctx context.Context) *UniBeeContext {
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*UniBeeContext); ok {
		return localCtx
	}
	return nil
}
