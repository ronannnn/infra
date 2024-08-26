package constant

type CtxKey string

const (
	CtxKeyAcceptLanguage CtxKey = "Accept-Language"
	CtxKeyLang           CtxKey = "Lang"
	CtxKeyUa             CtxKey = "User-Agent"
	CtxKeyDeviceId       CtxKey = "Device-Id"
)
