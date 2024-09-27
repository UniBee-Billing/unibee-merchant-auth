package jwt

import (
	"context"
	"fmt"
	"github.com/UniBee-Billing/unibee-merchant-auth/bean"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackyang-hk/go-tools/utility"
	"strings"
)

const (
	TOKEN_PREFIX            = "UniBee.Portal."
	TOKENTYPEUSER           = "USER"
	TOKENTYPEMERCHANTMember = "MERCHANT_MEMBER"
)

var Key = "3^&secret-key-for-UniBee*1!8*"
var Env = ""

type TokenType string

type TokenClaims struct {
	TokenType     TokenType `json:"tokenType"`
	Id            uint64    `json:"id"`
	Email         string    `json:"email"`
	MerchantId    uint64    `json:"merchantId"`
	PermissionKey string    `json:"permissionKey"`
	Lang          string    `json:"lang"`
	jwt.RegisteredClaims
}

func IsPortalToken(token string) bool {
	return strings.HasPrefix(token, TOKEN_PREFIX)
}

func ParsePortalToken(accessToken string) *TokenClaims {
	utility.Assert(len(Key) > 0, "server error: tokenKey is nil")
	accessToken = strings.Replace(accessToken, TOKEN_PREFIX, "", 1)
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Key), nil
	})
	return parsedAccessToken.Claims.(*TokenClaims)
}

func GetMemberPermissionKey(one *bean.MerchantMember) string {
	permissionKey := fmt.Sprintf("%v_%s", one.IsOwner, utility.MD5(utility.MarshalToJsonString(one.Permissions)))
	return permissionKey
}

func getAuthTokenRedisKey(token string) string {
	return fmt.Sprintf("auth#%s#%s", Env, token)
}

func IsAuthTokenAvailable(ctx context.Context, token string) bool {
	get, err := g.Redis().Get(ctx, getAuthTokenRedisKey(token))
	if err != nil {
		return false
	}
	if get != nil && len(get.String()) > 0 {
		return true
	}
	return false
}

func GetOpenApiKeyRedisKey(token string) string {
	return fmt.Sprintf("openApiKey#%s#%s", Env, token)
}
