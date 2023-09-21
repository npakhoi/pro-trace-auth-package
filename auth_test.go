package baseauth

import (
	"fmt"
	"testing"
)

var svc IAuthService

const url = "http://localhost:8080"

func init() {
	svc, _ = SetUpBaseAuthService(url)
}

func TestAuthService_Login(t *testing.T) {
	res := svc.Login("admin", "admin")
	fmt.Println("Login - res", res)
}

func TestAuthService_VerifyToken(t *testing.T) {

}

func TestAuthService_RefreshToken(t *testing.T) {

}
