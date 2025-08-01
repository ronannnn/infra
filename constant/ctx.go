package constant

type CtxKey string

const (
	CtxKeyAcceptLanguage CtxKey = "Accept-Language"
	CtxKeyUa             CtxKey = "User-Agent"
	CtxKeyDeviceId       CtxKey = "Device-Id"
	CtxKeyShowType       CtxKey = "Show-Type"
	CtxKeyUserId         CtxKey = "ctxKeyUserId"
	CtxKeyUsername       CtxKey = "ctxKeyUsername"
	CtxKeyClaims         CtxKey = "ctxKeyClaims"
)
