package model

import "github.com/dgrijalva/jwt-go"

type JwtCustomClaims struct {
	UserId             string `json:"user_id"`
	Role               string `json:"role"`
	jwt.StandardClaims        // token
}
