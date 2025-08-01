package handler

type RespCode int

const (
	NoErrorCode              RespCode = 0
	NormalErrorCode          RespCode = 1
	FieldValidationErrorCode RespCode = 2
	AccessTokenErrorCode     RespCode = 10
	RefreshTokenErrorCode    RespCode = 11
)

type ShowType int

const (
	ShowTypeSilent ShowType = 0

	ShowTypeSuccessMessage ShowType = 1
	ShowTypeInfoMessage    ShowType = 2
	ShowTypeWarningMessage ShowType = 3
	ShowTypeErrorMessage   ShowType = 4

	ShowTypeSuccessNotification ShowType = 11
	ShowTypeInfoNotification    ShowType = 12
	ShowTypeWarningNotification ShowType = 13
	ShowTypeErrorNotification   ShowType = 14
)

// Response According to https://pro.ant.design/zh-CN/docs/request
type Response struct {
	Success  bool     `json:"success"`
	Data     any      `json:"data,omitempty"`
	Message  string   `json:"message,omitempty"`
	Code     RespCode `json:"code"`
	ShowType ShowType `json:"showType"`
}
