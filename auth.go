package baseauth

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/npakhoi/pro-trace-auth-package/common"
	"github.com/npakhoi/pro-trace-auth-package/function/note"
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
	VerifyToken() gin.HandlerFunc
}

type authService struct {
	verifyTokenFunction verify.IVerifyFunction
}

func SetUpBaseAuthService(host string) (IAuthService, error) {
	err := initConnection(host)
	if err != nil {
		return nil, err
	}
	noteService := note.NewNoteFunction(host)
	verifyTokenFunction := verify.NewVerifyFunction(noteService, secretKey)
	return authService{
		verifyTokenFunction: verifyTokenFunction,
	}, nil
}

func (a authService) VerifyToken() gin.HandlerFunc {
	return a.verifyTokenFunction.VerifyToken()
}
