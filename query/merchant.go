package query

import (
	"context"
	"github.com/UniBee-Billing/unibee-merchant-auth/bean"
	"github.com/UniBee-Billing/unibee-merchant-auth/jwt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jackyang-hk/go-tools/utility"
	"strings"
)

func GetMerchantByApiKey(ctx context.Context, apiKey string) (one *bean.Merchant) {
	if len(apiKey) <= 0 {
		return nil
	}
	idData, err := g.Redis().Get(ctx, jwt.GetOpenApiKeyRedisKey(apiKey))
	data, err := g.Redis().Get(ctx, "UniBee#AllMerchants")
	if err != nil {
		return nil
	}
	var list []*bean.Merchant
	err = utility.UnmarshalFromJsonString(data.String(), &list)
	if err != nil {
		return nil
	}

	for _, merchant := range list {
		if idData != nil && idData.Uint64() > 0 {
			if merchant.Id == idData.Uint64() {
				one = merchant
				break
			}
		}
		if strings.Compare(merchant.ApiKey, apiKey) == 0 {
			one = merchant
			break
		}
	}
	return one
}

func GetMerchantById(ctx context.Context, id uint64) (one *bean.Merchant) {
	data, err := g.Redis().Get(ctx, "UniBee#AllMerchants")
	if err != nil {
		return nil
	}
	var list []*bean.Merchant
	err = utility.UnmarshalFromJsonString(data.String(), &list)
	if err != nil {
		return nil
	}

	for _, merchant := range list {
		if merchant.Id == id {
			one = merchant
			break
		}
	}
	return one
}
