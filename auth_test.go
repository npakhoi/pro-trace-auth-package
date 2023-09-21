package baseauth

import (
	"fmt"
	"testing"
)

var svc IAuthService
var token string
var refreshToken string

const url = "http://localhost:8080"

func init() {
	svc = SetUpBaseAuthService(url)
}

func TestAuthService_Login(t *testing.T) {
	res := svc.Login("admin", "admin")
	token = res.Token
	refreshToken = res.RefreshToken
	fmt.Println("Login - token", token)
	fmt.Println("Login - refreshToken", refreshToken)
}

func TestAuthService_VerifyToken(t *testing.T) {

}

func TestAuthService_RefreshToken(t *testing.T) {

}
