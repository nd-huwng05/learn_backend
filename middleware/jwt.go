package middleware

import (
	"backend-github-trending/model"
	"backend-github-trending/security"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func JWTMiddleware() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &model.JwtCustomClaims{},
		SigningKey: []byte(security.SECURITY_KEY),
	}
	return middleware.JWTWithConfig(config)
}
