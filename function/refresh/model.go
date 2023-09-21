package refresh

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenRefreshRq struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenRefreshRs struct {
	Token string `json:"token"`
}

type JwtClaims struct {
	ClaimId  uuid.UUID `json:"claimId,omitempty"`
	UserName string    `json:"userName,omitempty"`
	jwt.RegisteredClaims
}

type Response struct {
	Data   TokenRefreshRs `json:"data"`
	Status int            `json:"status"`
	Error  string         `json:"error"`
}
