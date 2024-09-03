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
	Silent ShowType = 0

	SuccessMessage ShowType = 1
	InfoMessage    ShowType = 2
	WarningMessage ShowType = 3
	ErrorMessage   ShowType = 4

	SuccessNotification ShowType = 11
	InfoNotification    ShowType = 12
	WarningNotification ShowType = 13
	ErrorNotification   ShowType = 14
)

// Response According to https://pro.ant.design/zh-CN/docs/request
type Response struct {
	Success  bool     `json:"success"`
	Data     any      `json:"data,omitempty"`
	Message  string   `json:"message,omitempty"`
	Code     RespCode `json:"code"`
	ShowType ShowType `json:"showType"`
}
