package jwt

import (
	"context"
	"fmt"
	"github.com/UniBee-Billing/unibee-merchant-auth"
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

var jwtKey = "3^&secret-key-for-UniBee*1!8*"
var env = ""

func SetupJwtToken(_env string, _jwtKey string) {
	jwtKey = _jwtKey
	env = _env
}

func IsPortalToken(token string) bool {
	return strings.HasPrefix(token, TOKEN_PREFIX)
}

func ParsePortalToken(accessToken string) *unibee_merchant_auth.TokenClaims {
	utility.Assert(len(jwtKey) > 0, "server error: tokenKey is nil")
	accessToken = strings.Replace(accessToken, TOKEN_PREFIX, "", 1)
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &unibee_merchant_auth.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	return parsedAccessToken.Claims.(*unibee_merchant_auth.TokenClaims)
}

func GetMemberPermissionKey(one *bean.MerchantMember) string {
	permissionKey := fmt.Sprintf("%v_%s", one.IsOwner, utility.MD5(utility.MarshalToJsonString(one.Permissions)))
	return permissionKey
}

func getAuthTokenRedisKey(token string) string {
	return fmt.Sprintf("auth#%s#%s", env, token)
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
	return fmt.Sprintf("openApiKey#%s#%s", env, token)
}
