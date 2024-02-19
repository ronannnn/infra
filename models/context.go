package models

type Key string

var (
	// context keys
	CtxKeyLang     Key = "ctxKeyLang"
	CtxKeyUserId   Key = "ctxKeyUserId"
	CtxKeyUsername Key = "ctxKeyUsername"
	CtxKeyClaims   Key = "ctxKeyClaims"
)
