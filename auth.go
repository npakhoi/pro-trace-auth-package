package baseauth

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/npakhoi/pro-trace-auth-package/common"
	"github.com/npakhoi/pro-trace-auth-package/function/login"
	"github.com/npakhoi/pro-trace-auth-package/function/note"
	"github.com/npakhoi/pro-trace-auth-package/function/refresh"
	"github.com/npakhoi/pro-trace-auth-package/function/verify"
	"io"
	"net/http"
)

var secretKey string

func initConnection(host string) error {
	url := host + "/get-secret"
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check the response status code.
	if response.StatusCode != http.StatusOK {
		return errors.New("failed")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var res common.Response

	if err = json.Unmarshal(body, &res); err != nil {
		return err
	} else {
		secretKey = res.Data
	}
	return nil
}

type IAuthService interface {
	Login(username string, password string) login.Response
	VerifyToken() gin.HandlerFunc
	RefreshToken(rfToken string) refresh.TokenRefreshRs
}

type authService struct {
	loginService        login.ILoginFunction
	verifyTokenFunction verify.IVerifyFunction
	refreshService      refresh.IRefreshFunction
}

func SetUpBaseAuthService(host string) (IAuthService, error) {
	err := initConnection(host)
	if err != nil {
		return nil, err
	}
	loginService := login.NewLoginService(host)
	refreshService := refresh.NewRefreshService(host)
	noteService := note.NewNoteFunction(host)
	verifyTokenFunction := verify.NewVerifyFunction(noteService, secretKey)
	return authService{
		loginService:        loginService,
		verifyTokenFunction: verifyTokenFunction,
		refreshService:      refreshService,
	}, nil
}

func (a authService) Login(username string, password string) login.Response {
	cred := login.Credentials{
		UserName: username,
		Password: password,
	}
	return a.loginService.Login(cred)
}

func (a authService) VerifyToken() gin.HandlerFunc {
	return a.verifyTokenFunction.VerifyToken()
}

func (a authService) RefreshToken(rfToken string) refresh.TokenRefreshRs {
	return a.refreshService.ResetToken(rfToken)
}
