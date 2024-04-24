package models

type Key string

var (
	// context keys
	CtxKeyLang     Key = "ctxKeyLang"
	CtxKeyUa       Key = "ctxKeyUa"
	CtxKeyDeviceId Key = "ctxKeyDeviceId"
	CtxKeyUserId   Key = "ctxKeyUserId"
	CtxKeyUsername Key = "ctxKeyUsername"
	CtxKeyClaims   Key = "ctxKeyClaims"
)
