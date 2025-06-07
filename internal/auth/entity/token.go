package entity

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}
