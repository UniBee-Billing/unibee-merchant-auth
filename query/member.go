package query

import (
	"context"
	"fmt"
	"github.com/UniBee-Billing/unibee-merchant-auth/bean"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jackyang-hk/go-tools/utility"
)

func GetMerchantMemberById(ctx context.Context, id uint64) (one *bean.MerchantMember) {
	if id <= 0 {
		return nil
	}
	data, err := g.Redis().Get(ctx, fmt.Sprintf("UniBee#Member#%d", id))
	if err != nil {
		return nil
	}
	err = utility.UnmarshalFromJsonString(data.String(), &one)
	if err != nil {
		return nil
	}

	return one
}
