package token

import "github.com/golang-jwt/jwt/v5"


type UserClaims struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}