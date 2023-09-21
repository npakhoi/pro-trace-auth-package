package verify

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshRq struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshRs struct {
	Token string `json:"token"`
}

type JwtClaims struct {
	ClaimId  uuid.UUID `json:"claimId,omitempty"`
	UserName string    `json:"userName,omitempty"`
	jwt.RegisteredClaims
}
