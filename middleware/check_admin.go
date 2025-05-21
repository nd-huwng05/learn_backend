package middleware

import (
	"backend-github-trending/log"
	"backend-github-trending/model"
	req2 "backend-github-trending/model/req"
	"github.com/labstack/echo/v4"
	"net/http"
)

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := req2.ReqSignIn{}
			if err := c.Bind(&req); err != nil {
				log.Log.Error(err.Error())
				return c.JSON(http.StatusBadRequest, model.Response{
					StausCode: http.StatusBadRequest,
					Message:   err.Error(),
					Data:      nil,
				})
			}

			if req.Email != "admin@gmail.com" {
				return c.JSON(http.StatusBadRequest, model.Response{
					StausCode: http.StatusBadRequest,
					Message:   "You isn't admin, you don't have permission to access this page",
					Data:      nil,
				})
			}

			return next(c)
		}
	}
}
