package login

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
	"github.com/ronannnn/infra/service/jwt"
	"github.com/ronannnn/infra/service/jwt/refreshtoken"
	"github.com/ronannnn/infra/utils"
	"github.com/ronannnn/infra/utils/useragent"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewService(
	log *zap.SugaredLogger,
	jwtService *jwt.Service,
	loginRecordService *service.LoginRecordService,
	userService *service.UserService,
) *Service {
	return &Service{
		log:                log,
		jwtService:         jwtService,
		loginRecordService: loginRecordService,
		userService:        userService,
	}
}

type Service struct {
	log                *zap.SugaredLogger
	jwtService         *jwt.Service
	loginRecordService *service.LoginRecordService
	userService        *service.UserService
}

func (srv *Service) LoginByUsernameAndPassword(r *http.Request, username, password string) (loginResp *Result, err error) {
	ctx := r.Context()

	var user *model.User
	if user, err = srv.userService.GetByUsername(ctx, username); err == gorm.ErrRecordNotFound {
		return nil, model.ErrWrongUsernameOrPassword
	} else if err != nil {
		return
	}

	// check if login type is allowed for this user
	if err = user.HasLoginType(LOGIN_TYPE_USERNAME_PASSWORD); err != nil {
		return
	}

	// get login info
	ua, deviceId, loginRecord := srv.getLoginInfo(r, LOGIN_TYPE_USERNAME_PASSWORD, user.Id)
	defer func() {
		if createdErr := srv.loginRecordService.Create(ctx, loginRecord); createdErr != nil {
			srv.log.Warnf("failed to create login record: %v", err)
		}
	}()

	// login validation
	if user.Password == nil {
		err = fmt.Errorf("密码未设置")
		return
	}
	if !utils.CheckPassword(*user.Password, password) {
		err = fmt.Errorf("密码不正确")
		return
	}

	// generate token
	var refreshToken, accessToken string
	var dupLogin bool
	if refreshToken, accessToken, dupLogin, err = srv.jwtService.GenerateTokens(ctx, refreshtoken.BaseClaims{
		UserId:   user.Id,
		Username: *user.Username,
	}, ua, deviceId); err != nil {
		return
	}
	if dupLogin {
		status := model.LoginStatusDupLogin
		loginRecord.Status = &status
	} else {
		status := model.LoginStatusSuccess
		loginRecord.Status = &status
	}
	loginResp = &Result{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		DupLogin:     dupLogin,
		User:         user,
	}
	return
}

func (srv *Service) getLoginInfo(r *http.Request, loginType string, userId uint) (ua, deviceId string, loginRecord *model.LoginRecord) {
	ctx := r.Context()
	ua = ctx.Value(constant.CtxKeyUa).(string)
	deviceId = ctx.Value(constant.CtxKeyDeviceId).(string)
	loginDeviceType := useragent.Parse(ua).DeviceType()
	loginTime := time.Now()
	status := model.LoginStatusFailed
	loginRecord = &model.LoginRecord{
		Ip:              &r.RemoteAddr,
		LoginTime:       &loginTime,
		DeviceId:        &deviceId,
		UserAgent:       &ua,
		LoginDeviceType: &loginDeviceType,
		Status:          &status,
		LoginType:       &loginType,
		UserId:          &userId,
	}
	return
}

func (srv *Service) Logout(ctx context.Context, userId uint, userAgent string) (err error) {
	return srv.jwtService.DeleteTokenByUserIdAndLoginDeviceType(ctx, userId, useragent.Parse(userAgent).DeviceType())
}

func (srv *Service) ChangePwd(ctx context.Context, cmd ChangeUserPwdCmd) (err error) {
	var user *model.User
	if user, err = srv.userService.GetById(ctx, cmd.UserId); err != nil {
		return
	}
	if !CheckPassword(*user.Password, cmd.OldPwd) {
		return model.ErrWrongUsernameOrPassword
	}
	var hashedNewPwd string
	if hashedNewPwd, err = HashPassword(cmd.NewPwd); err != nil {
		return
	}
	return srv.userService.ChangePwd(ctx, cmd.UserId, hashedNewPwd)
}
