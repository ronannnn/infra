package response

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	NoErrorCode           = 0
	NormalErrorCode       = 1
	AccessTokenErrorCode  = 10
	RefreshTokenErrorCode = 11
)

type ShowType int

const (
	Silent              ShowType = 0
	WarningMessage      ShowType = 1
	ErrorMessage        ShowType = 2
	WarningNotification ShowType = 4
	ErrorNotification   ShowType = 5
	Redirect            ShowType = 9
)

// Response According to https://pro.ant.design/zh-CN/docs/request
type Response struct {
	Success  bool     `json:"success"`
	Data     any      `json:"data,omitempty"`
	Message  string   `json:"message,omitempty"`
	Code     int      `json:"code,omitempty"`
	ShowType ShowType `json:"showType"`
}

func Result(w http.ResponseWriter, r *http.Request, response Response) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

func Ok(w http.ResponseWriter, r *http.Request) {
	Result(w, r, Response{Success: true})
}

func OkWithWarningNotifMsg(w http.ResponseWriter, r *http.Request, msg string) {
	Result(w, r, Response{Success: true, Message: msg, ShowType: WarningNotification})
}

func OkWithErrorNotifMsg(w http.ResponseWriter, r *http.Request, msg string) {
	Result(w, r, Response{Success: true, Message: msg, ShowType: ErrorNotification})
}

func OkWithDataAndWarningNotifMsg(w http.ResponseWriter, r *http.Request, data any, msg string) {
	Result(w, r, Response{Success: true, Data: data, Message: msg, ShowType: WarningNotification})
}

func OkWithDataAndErrorNotifMsg(w http.ResponseWriter, r *http.Request, data any, msg string) {
	Result(w, r, Response{Success: true, Data: data, Message: msg, ShowType: ErrorNotification})
}

func OkWithData(w http.ResponseWriter, r *http.Request, data any) {
	Result(w, r, Response{Success: true, Data: data})
}

func FailWithMsg(w http.ResponseWriter, r *http.Request, msg string) {
	Result(w, r, Response{Message: msg, Code: NormalErrorCode, ShowType: ErrorMessage})
}

func FailWithErr(w http.ResponseWriter, r *http.Request, err error) {
	Result(w, r, Response{Message: err.Error(), Code: NormalErrorCode, ShowType: ErrorMessage})
}

func FailWithSilentErr(w http.ResponseWriter, r *http.Request, err error) {
	Result(w, r, Response{Message: err.Error(), Code: NormalErrorCode, ShowType: Silent})
}

func ErrAccessToken(w http.ResponseWriter, r *http.Request, err error) {
	Result(w, r, Response{Message: err.Error(), Code: AccessTokenErrorCode, ShowType: ErrorMessage})
}

func ErrRefreshToken(w http.ResponseWriter, r *http.Request, err error) {
	Result(w, r, Response{Message: err.Error(), Code: RefreshTokenErrorCode, ShowType: ErrorMessage})
}

func ErrPrivilege(w http.ResponseWriter, r *http.Request) {
	Result(w, r, Response{Message: "权限不足", Code: NormalErrorCode, ShowType: ErrorMessage})
}
