package verify

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/npakhoi/pro-trace-auth-package/function/note"
	"net/http"
	"strings"
)

type IVerifyFunction interface {
	VerifyToken() gin.HandlerFunc
}

type function struct {
	note      note.INoteFunction
	secretKey string
}

func NewVerifyFunction(note note.INoteFunction, secretKey string) IVerifyFunction {
	return function{note, secretKey}
}

type authHeader struct {
	Authorization string `header:"Authorization"`
}

func (f function) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}
		if err := c.ShouldBindHeader(&h); err != nil {
			MapErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		authParts := strings.Split(h.Authorization, " ")

		if h.Authorization == "" || len(authParts) < 2 || authParts[0] == "" || authParts[1] == "" {
			MapErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("token is required"))
			return
		}

		if !strings.HasPrefix(h.Authorization, "Bearer ") {
			MapErrorResponse(c, http.StatusBadRequest, jwt.ErrTokenMalformed)
			return
		}

		claims := &JwtClaims{}
		parsedToken, err := jwt.ParseWithClaims(authParts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(f.secretKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				MapErrorResponse(c, http.StatusUnauthorized, err)
				return
			}
			if err == jwt.ErrTokenExpired || err.Error() == "token has invalid claims: token is expired" {
				f.note.NoteExpiredToken(authParts[1])
				MapErrorResponse(c, http.StatusUnauthorized, errors.New("token is expired"))
				return
			} else {
				MapErrorResponse(c, http.StatusBadRequest, err)
				return
			}
		}
		if !parsedToken.Valid {
			MapErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("access token is invalid"))
			return
		}
		c.Next()
	}
}

type Response struct {
	Data   any    `json:"data"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func MapErrorResponse(ctx *gin.Context, code int, err error) {
	response := Response{
		Data:   "",
		Status: code,
		Error:  err.Error(),
	}
	ctx.JSON(code, response)
	ctx.Abort()
}
