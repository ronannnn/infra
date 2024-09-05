package constant

type CtxKey string

const (
	CtxKeyAcceptLanguage CtxKey = "Accept-Language"
	CtxKeyUa             CtxKey = "User-Agent"
	CtxKeyDeviceId       CtxKey = "Device-Id"
	CtxKeyUserId         CtxKey = "ctxKeyUserId"
	CtxKeyUsername       CtxKey = "ctxKeyUsername"
	CtxKeyClaims         CtxKey = "ctxKeyClaims"
)
