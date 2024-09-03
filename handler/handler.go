package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/i18n"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/validator"
	"go.uber.org/zap"
)

type HttpHandler interface {
	// handle request binding and checking
	BindAndCheck(w http.ResponseWriter, r *http.Request, data any) bool
	BindUint64Param(w http.ResponseWriter, r *http.Request, key string, data *uint64) bool
	BindParam(w http.ResponseWriter, r *http.Request, key string, data *string) bool

	// handle response
	Success(w http.ResponseWriter, r *http.Request, message *msg.Message, data any)
	SuccessSilently(w http.ResponseWriter, r *http.Request, message *msg.Message, data any)
	SuccessWithShowType(w http.ResponseWriter, r *http.Request, message *msg.Message, data any, showType ShowType)
	// err may be msg.Error, which has *msg.Message
	Fail(w http.ResponseWriter, r *http.Request, err error, data any)
	FailWithCode(w http.ResponseWriter, r *http.Request, err error, data any, code RespCode)
	FailWithShowType(w http.ResponseWriter, r *http.Request, err error, data any, showType ShowType)
	FailWithCodeAndShowType(w http.ResponseWriter, r *http.Request, err error, data any, code RespCode, showType ShowType)
}

func NewHttpHandler(
	log *zap.SugaredLogger,
	i18n i18n.I18n,
	validator validator.Validator,
) HttpHandler {
	return &HttpHandlerImpl{
		log:       log,
		i18n:      i18n,
		validator: validator,
	}
}

type HttpHandlerImpl struct {
	log       *zap.SugaredLogger
	i18n      i18n.I18n
	validator validator.Validator
}

// BindAndCheck bind request and check
func (h *HttpHandlerImpl) BindAndCheck(w http.ResponseWriter, r *http.Request, data any) bool {
	var err error
	lang := GetLang(r)
	r = r.WithContext(context.WithValue(r.Context(), constant.CtxKeyAcceptLanguage, lang))
	if err = render.DefaultDecoder(r, &data); err != nil {
		h.Fail(w, r, msg.NewError(reason.RequestFormatError), err)
		return true
	}
	var errFields []*validator.FormErrorField
	if errFields, err = h.validator.Check(r.Context(), lang, data); err != nil {
		h.Fail(w, r, err, errFields)
		return true
	}
	return false
}

// BindAndCheck bind request and check
func (h *HttpHandlerImpl) BindAndCheckPartial(w http.ResponseWriter, r *http.Request, data any) bool {
	var err error
	lang := GetLang(r)
	r = r.WithContext(context.WithValue(r.Context(), constant.CtxKeyAcceptLanguage, lang))
	if err = render.DefaultDecoder(r, &data); err != nil {
		h.Fail(w, r, msg.NewError(reason.RequestFormatError), err)
		return true
	}
	var errFields []*validator.FormErrorField
	if errFields, err = h.validator.CheckPartial(r.Context(), lang, data); err != nil {
		h.Fail(w, r, err, errFields)
		return true
	}
	return false
}

func (h *HttpHandlerImpl) BindUint64Param(w http.ResponseWriter, r *http.Request, key string, data *uint64) bool {
	param := chi.URLParam(r, key)
	if param == "" {
		h.Fail(w, r, msg.NewError(reason.MissingRequiredParam).WithReasonTemplateData(reason.ParamTd{Param: key}).WithStack(), nil)
		return true
	}
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		h.Fail(w, r, msg.NewError(reason.InvalidUintParam).WithReasonTemplateData(reason.ParamTd{Param: key}).WithStack(), nil)
		return true
	}
	*data = id
	return false
}

func (h *HttpHandlerImpl) BindParam(w http.ResponseWriter, r *http.Request, key string, data *string) bool {
	param := chi.URLParam(r, key)
	if param == "" {
		h.Fail(w, r, msg.NewError(reason.MissingRequiredParam).WithReasonTemplateData(reason.ParamTd{Param: key}).WithStack(), nil)
		return true
	}
	*data = param
	return false
}

func (h *HttpHandlerImpl) Success(
	w http.ResponseWriter,
	r *http.Request,
	message *msg.Message,
	data any,
) {
	h.SuccessWithShowType(w, r, message, data, SuccessMessage)
}

func (h *HttpHandlerImpl) SuccessSilently(
	w http.ResponseWriter,
	r *http.Request,
	message *msg.Message,
	data any,
) {
	h.SuccessWithShowType(w, r, message, data, Silent)
}

func (h *HttpHandlerImpl) SuccessWithShowType(
	w http.ResponseWriter,
	r *http.Request,
	message *msg.Message,
	data any,
	showType ShowType,
) {
	respMsg := getRespMsg(r, h.i18n, message)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Response{Success: true, Data: data, Message: respMsg, ShowType: showType})
}

// Fail handle response
func (h *HttpHandlerImpl) Fail(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	data any,
) {
	h.FailWithCodeAndShowType(w, r, err, data, NormalErrorCode, ErrorMessage)
}

// Fail handle response
func (h *HttpHandlerImpl) FailWithCode(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	data any,
	code RespCode,
) {
	h.FailWithCodeAndShowType(w, r, err, data, code, ErrorNotification)
}

// Fail handle response
func (h *HttpHandlerImpl) FailWithShowType(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	data any,
	showType ShowType,
) {
	h.FailWithCodeAndShowType(w, r, err, data, NormalErrorCode, showType)
}

// Fail handle response
func (h *HttpHandlerImpl) FailWithCodeAndShowType(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	data any,
	code RespCode,
	showType ShowType,
) {

	// 有错误，返回错误信息
	var msgErr *msg.Error
	// unknown error
	if !errors.As(err, &msgErr) {
		h.log.Error(err, "\n", msg.LogStack(2, 5))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{Message: err.Error(), Code: code, ShowType: showType})
		return
	}

	// known error
	if msgErr.Err != nil {
		h.log.Error(msgErr.Err, "\n", msg.LogStack(2, 3))
	}
	if errors.Is(msgErr.Err, constant.ErrFieldsValidation) {
		// validation error
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{Code: FieldValidationErrorCode, Data: data})
		return
	}
	respMsg := getRespMsg(r, h.i18n, &msgErr.Message)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Response{Message: respMsg, Code: code, ShowType: showType})
}

func getRespMsg(r *http.Request, i18n i18n.I18n, message *msg.Message) string {
	if message == nil {
		return ""
	}
	lang := GetLang(r)
	if message.Msg != "" {
		return message.Msg
	}
	return i18n.TrWithData(lang, message.Reason, message.ReasonTemplateData)
}
