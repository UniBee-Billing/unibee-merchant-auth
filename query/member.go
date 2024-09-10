package query

import (
	"context"
	"github.com/UniBee-Billing/unibee-merchant-auth/bean"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jackyang-hk/go-tools/utility"
)

func GetMerchantMemberById(ctx context.Context, id uint64) (one *bean.MerchantMember) {
	if id <= 0 {
		return nil
	}
	data, err := g.Redis().Get(ctx, "UniBee#AllMembers")
	if err != nil {
		return nil
	}
	var list []*bean.MerchantMember
	err = utility.UnmarshalFromJsonString(data.String(), &list)
	if err != nil {
		return nil
	}

	for _, member := range list {
		if member.Id == id {
			one = member
			break
		}
	}

	return one
}
