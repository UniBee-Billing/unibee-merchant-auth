package unibee_merchant_auth

import (
	"context"
	"fmt"
	"github.com/UniBee-Billing/unibee-merchant-auth/i18n"
	"github.com/UniBee-Billing/unibee-merchant-auth/jwt"
	"github.com/UniBee-Billing/unibee-merchant-auth/query"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	_ "github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jackyang-hk/go-tools/utility"
	"strconv"
	"strings"
)

var loginUrl = "/login"

func CORS(r *ghttp.Request) {
	g.Log().Debugf(r.Context(), "CORS Control: HTTP Header Host:%s", r.GetHost())
	g.Log().Debugf(r.Context(), "CORS Control: HTTP Header Origin:%s", r.GetHeader("Origin"))
	g.Log().Debugf(r.Context(), "CORS Control: HTTP Header Referer:%s", r.GetHeader("Referer"))
	g.Log().Debugf(r.Context(), "CORS Control: HTTP Header User-Agent:%s", r.GetHeader("User-Agent"))
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func ResponseHandler(r *ghttp.Request) {
	customCtx := &UniBeeContext{
		Session: r.Session,
		Data:    make(g.Map),
	}
	customCtx.RequestId = utility.CreateRequestId()
	Context().Init(r, customCtx)
	r.Assigns(g.Map{
		ContextKey: customCtx,
	})

	// Setup System Default Language
	r.SetCtx(gi18n.WithLanguage(r.Context(), "en"))
	customCtx.Language = "en"
	lang := ""
	if r.Get("lang") != nil {
		lang = r.Get("lang").String()
	}
	if len(lang) == 0 {
		lang = r.GetHeader("lang")
	}
	if len(lang) > 0 && i18n.IsLangAvailable(lang) {
		r.SetCtx(gi18n.WithLanguage(r.Context(), strings.ToLower(strings.TrimSpace(lang))))
		customCtx.Language = lang
	}

	customCtx.UserAgent = r.Header.Get("User-Agent")
	if len(customCtx.UserAgent) > 0 && strings.Contains(customCtx.UserAgent, "OpenAPI") {
		customCtx.IsOpenApiCall = true
	}
	customCtx.Authorization = r.Header.Get("Authorization")
	customCtx.TokenString = customCtx.Authorization
	if len(customCtx.TokenString) > 0 && strings.HasPrefix(customCtx.TokenString, "Bearer ") && !jwt.IsPortalToken(customCtx.TokenString) {
		customCtx.IsOpenApiCall = true
		customCtx.TokenString = strings.Replace(customCtx.TokenString, "Bearer ", "", 1) // remove Bearer
	}
	g.Log().Info(r.Context(), fmt.Sprintf("[Request][%s][%s][%s][%s] IsOpenApi:%v Token:%s Body:%s", customCtx.Language, customCtx.RequestId, r.Method, r.GetUrl(), customCtx.IsOpenApiCall, customCtx.TokenString, r.GetBodyString()))

	utility.Try(r.Middleware.Next, func(err interface{}) {
		g.Log().Errorf(r.Context(), "[Request][%s][%s][%s] Global_Exception Panic Body:%s Error:%v", customCtx.RequestId, r.Method, r.GetUrl(), r.GetBodyString(), err)
		return
	})
	g.Log().Info(r.Context(), fmt.Sprintf("[Request][%s][%s][%s] MerchantId:%d", customCtx.RequestId, r.Method, r.GetUrl(), customCtx.MerchantId))

	var (
		err             = r.GetError()
		res             = r.GetHandlerResponse()
		code gcode.Code = gcode.CodeOK
	)

	if err == nil && r.Response.BufferLength() > 0 {
		return
	}

	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		json, _ := r.GetJson()
		g.Log().Errorf(r.Context(), "Global_exception requestId:%s url: %s params:%s code:%d error:%s", Context().GetUniBeeContext(r.Context()).RequestId, r.GetUrl(), json, code.Code(), err.Error())
		r.Response.ClearBuffer() // inner panic will contain json data，need clean

		message := err.Error()
		if strings.Contains(message, "Session Expired") {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				OpenApiJsonExit(r, gcode.CodeValidationFailed.Code(), "Session Expired")
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				JsonRedirectExit(r, 61, "Session Expired", loginUrl)
			}
		} else if strings.Contains(message, utility.SystemAssertPrefix) || code == gcode.CodeValidationFailed {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				OpenApiJsonExit(r, gcode.CodeValidationFailed.Code(), strings.Replace(message, "exception recovered: "+utility.SystemAssertPrefix, "", 1))
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				JsonExit(r, gcode.CodeValidationFailed.Code(), strings.Replace(message, "exception recovered: "+utility.SystemAssertPrefix, "", 1))
			}
		} else {
			if customCtx.IsOpenApiCall {
				r.Response.Status = 400
				OpenApiJsonExit(r, code.Code(), fmt.Sprintf("Server Error-%s-%d", Get(r.Context()).RequestId, code.Code()))
			} else {
				r.Response.Status = 200 // error reply in json code, http code always 200
				JsonExit(r, code.Code(), fmt.Sprintf("Server Error-%s-%d", Get(r.Context()).RequestId, code.Code()))
			}
		}
	} else {
		r.Response.Status = 200
		if customCtx.IsOpenApiCall {
			OpenApiJsonExit(r, code.Code(), "", res)
		} else {
			JsonExit(r, code.Code(), "", res)
		}
	}
}

func MerchantHandler(r *ghttp.Request) {
	customCtx := Context().GetUniBeeContext(r.Context())
	if len(customCtx.TokenString) == 0 {
		g.Log().Infof(r.Context(), "MerchantHandler empty token string of auth header")
		if customCtx.IsOpenApiCall {
			r.Response.Status = 401
			OpenApiJsonExit(r, 61, "invalid token")
		} else {
			JsonRedirectExit(r, 61, "invalid token", loginUrl)
		}
		r.Exit()
	}
	if !customCtx.IsOpenApiCall {
		// Merchant Portal Call
		if !jwt.IsAuthTokenAvailable(r.Context(), customCtx.TokenString) {
			g.Log().Infof(r.Context(), "MerchantHandler Invalid Token:%s", customCtx.TokenString)
			JsonRedirectExit(r, 61, "invalid token", loginUrl)
			r.Exit()
		}

		customCtx.Token = jwt.ParsePortalToken(customCtx.TokenString)
		g.Log().Debugf(r.Context(), "MerchantHandler Parsed Token: %s, URL: %s", utility.MarshalToJsonString(customCtx.Token), r.GetUrl())

		if customCtx.Token.TokenType == jwt.TOKENTYPEMERCHANTMember {
			member := query.GetMerchantMemberById(r.Context(), customCtx.Token.Id)
			permissionKey := jwt.GetMemberPermissionKey(member)
			if member == nil {
				g.Log().Infof(r.Context(), "MerchantHandler merchant member not found token:%s", utility.MarshalToJsonString(customCtx.Token))
				JsonRedirectExit(r, 61, "merchant user not found", loginUrl)
				r.Exit()
			} else if member.Status == 2 {
				g.Log().Infof(r.Context(), "MerchantHandler merchant member has suspend :%v", utility.MarshalToJsonString(customCtx.Token))
				JsonRedirectExit(r, 61, "Your account has been suspended. Please contact billing admin for further assistance.", loginUrl)
				r.Exit()
			} else if strings.Compare(permissionKey, customCtx.Token.PermissionKey) != 0 && !strings.Contains(r.GetUrl(), "logout") {
				g.Log().Infof(r.Context(), "MerchantHandler merchant member permission has change, need reLogin")
				JsonRedirectExit(r, 62, "Your permission has changed. Please reLogin.", loginUrl)
				r.Exit()
			}

			customCtx.MerchantMember = &UniBeeContextMerchantMember{
				Id:         customCtx.Token.Id,
				MerchantId: customCtx.Token.MerchantId,
				Token:      customCtx.TokenString,
				Email:      customCtx.Token.Email,
				IsOwner:    strings.Compare(strings.Trim(member.Role, " "), "Owner") == 0,
			}
			customCtx.MerchantId = customCtx.Token.MerchantId
			doubleRequestLimit(strconv.FormatUint(customCtx.MerchantMember.Id, 10), r)
			lang := ""
			if r.Get("lang") != nil {
				lang = r.Get("lang").String()
			}
			if len(lang) == 0 {
				lang = r.GetHeader("lang")
			}
			if len(lang) > 0 && i18n.IsLangAvailable(lang) {
				r.SetCtx(gi18n.WithLanguage(r.Context(), strings.ToLower(strings.TrimSpace(lang))))
			}
		} else {
			g.Log().Infof(r.Context(), "MerchantHandler invalid token type token:%v", utility.MarshalToJsonString(customCtx.Token))
			JsonRedirectExit(r, 61, "invalid token type", loginUrl)
			r.Exit()
		}
	} else {
		// Api Call
		customCtx.IsOpenApiCall = true
		merchantInfo := query.GetMerchantByApiKey(r.Context(), customCtx.TokenString)
		if merchantInfo == nil {
			r.Response.Status = 401
			OpenApiJsonExit(r, 61, "invalid token")
		} else {
			customCtx.MerchantId = merchantInfo.Id
			customCtx.OpenApiKey = customCtx.TokenString
		}
		lang := ""
		if r.Get("lang") != nil {
			lang = r.Get("lang").String()
		}
		if len(lang) == 0 {
			lang = r.GetHeader("lang")
		}
		if len(lang) > 0 && i18n.IsLangAvailable(lang) {
			r.SetCtx(gi18n.WithLanguage(r.Context(), strings.ToLower(strings.TrimSpace(lang))))
		}
	}
	r.Middleware.Next()
}

func doubleRequestLimit(id string, r *ghttp.Request) {
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
		if strings.HasSuffix(r.GetUrl(), "detail") || strings.HasSuffix(r.GetUrl(), "list") || strings.HasSuffix(r.GetUrl(), "get") {
			return
		}
		md5 := utility.MD5(fmt.Sprintf("%s%s%s", id, r.GetUrl(), r.GetBodyString()))
		if !utility.TryLock(r.Context(), md5, 2) {
			utility.Assert(false, i18n.LocalizationFormat(r.Context(), "{#ClickTooFast}"))
		}
	}
}

func Get(ctx context.Context) *UniBeeContext {
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*UniBeeContext); ok {
		return localCtx
	}
	return nil
}
