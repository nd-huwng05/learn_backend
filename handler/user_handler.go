package handler

import (
	"backend-github-trending/banana"
	"backend-github-trending/log"
	"backend-github-trending/model"
	req2 "backend-github-trending/model/req"
	"backend-github-trending/repository"
	"backend-github-trending/security"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/google/uuid" // create uuid
	"github.com/labstack/echo"
	"net/http"
)

type UserHandler struct {
	UserRepo repository.UserRepo
}

func (u *UserHandler) HandleSignIn(c echo.Context) error {
	req := req2.ReqSignIn{}
	if err := c.Bind(&req); err != nil {
		log.Log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StausCode: http.StatusBadRequest,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	if err := c.Validate(req); err != nil {
		log.Log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StausCode: http.StatusBadRequest,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	user, err := u.UserRepo.CheckLogin(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StausCode: http.StatusUnauthorized,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	// check pass
	isTheSame := security.ComparePasswords(user.Password, []byte(req.Password))
	if !isTheSame {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StausCode: http.StatusUnauthorized,
			Message:   "Login Failed",
			Data:      nil,
		})
	}

	// gen token
	token, err := security.GenToken(user)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StausCode: http.StatusUnauthorized,
			Message:   err.Error(),
			Data:      nil,
		})
	}
	user.Token = token
	user.Password = ""

	return c.JSON(http.StatusOK, model.Response{
		StausCode: http.StatusOK,
		Message:   "Success",
		Data:      user,
	})
}

func (u *UserHandler) HandleSignUp(c echo.Context) error {
	req := req2.ReqSignUp{}
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StausCode: http.StatusBadRequest,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	if err := c.Validate(req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StausCode: http.StatusBadRequest,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	hash := security.HashAndSalt([]byte(req.Password)) // ma hoa du lieu
	role := model.MEMBER.String()

	userId, err := uuid.NewUUID() // tao ra uuid
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusForbidden, model.Response{
			StausCode: http.StatusForbidden,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	user := model.User{
		UserId:   userId.String(),
		FullName: req.FullName,
		Email:    req.Email,
		Password: hash,
		Role:     role,
		Token:    "",
	}

	user, err = u.UserRepo.SaveUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusConflict, model.Response{
			StausCode: http.StatusConflict,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	type User struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
	}

	// gen token
	token, err := security.GenToken(user)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StausCode: http.StatusUnauthorized,
			Message:   err.Error(),
			Data:      nil,
		})
	}
	user.Token = token
	user.Password = ""

	return c.JSON(http.StatusOK, model.Response{
		StausCode: http.StatusOK,
		Data:      user,
		Message:   "User Created",
	})
}

func (u *UserHandler) HandleProfile(c echo.Context) error {
	tokenData := c.Get("user").(*jwt.Token)
	claims := tokenData.Claims.(*model.JwtCustomClaims)

	user, err := u.UserRepo.SelectUserById(c.Request().Context(), claims.UserId)
	if err != nil {
		if err == banana.UserNotFound {
			return c.JSON(http.StatusUnauthorized, model.Response{
				StausCode: http.StatusUnauthorized,
				Message:   err.Error(),
				Data:      nil,
			})
		}

		return c.JSON(http.StatusInternalServerError, model.Response{
			StausCode: http.StatusInternalServerError,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	user.Password = ""
	return c.JSON(http.StatusOK, model.Response{
		StausCode: http.StatusOK,
		Message:   "Success",
		Data:      user,
	})
}

func (u *UserHandler) HandleUpdateProfile(c echo.Context) error {
	req := req2.ReqUpdateUser{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	err := c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StausCode: http.StatusBadRequest,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaims)
	user := model.User{
		UserId:   claims.UserId,
		FullName: req.FullName,
		Email:    req.Email,
	}

	user, err = u.UserRepo.UpdateUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StausCode: http.StatusUnprocessableEntity,
			Message:   err.Error(),
			Data:      nil,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StausCode: http.StatusOK,
		Message:   "Success",
		Data:      user,
	})
}
