## UniBee Backend Authentication Middleware Project

Based on github.com/gogf/gf
- go get "github.com/gogf/gf/v2/os/gctx"
- Doc: https://goframe.org/display/gf

## How to use

### Go Get:
``` 
    go get github.com/UniBee-Billing/unibee-merchant-auth@v1.0.7
``` 

### Setup:
``` 
unibee_merchant_auth.Setup("${YOUR Env}", "${YOUR TOKEN KEY}")
``` 

### Middleware Injection
```
group.Middleware(
					unibee_merchant_auth.CORS,
					unibee_merchant_auth.ResponseHandler,
					unibee_merchant_auth.MerchantHandler,
				)
```

### Get Merchant|Member|OpenApi Info
```
	g.Log().Infof(ctx, "merchantId:%d", unibee_merchant_auth.GetUniBeeContext(ctx).MerchantId)
	g.Log().Infof(ctx, "merchant:%s", utility.MarshalToJsonString(unibee_merchant_auth.GetUniBeeContext(ctx).Merchant))
	g.Log().Infof(ctx, "merchantMember:%d", utility.MarshalToJsonString(unibee_merchant_auth.GetUniBeeContext(ctx).MerchantMember))
```