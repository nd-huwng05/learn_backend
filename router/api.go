package router

import (
	"backend-github-trending/handler"
	"backend-github-trending/middleware"
	"github.com/labstack/echo"
)

type API struct {
	Echo        *echo.Echo
	UserHandler handler.UserHandler
}

func (api *API) SetupRoutes() {
	api.Echo.POST("/user/sign-in", api.UserHandler.HandleSignIn)
	api.Echo.POST("/user/sign-up", api.UserHandler.HandleSignUp)

	user := api.Echo.Group("/user", middleware.JWTMiddleware())
	user.POST("/profile", api.UserHandler.HandleProfile)
	user.POST("/profile/update", api.UserHandler.HandleUpdateProfile)
}
