package main

import (
	"backend-github-trending/db"
	"backend-github-trending/handler"
	"backend-github-trending/helper"
	"backend-github-trending/repository/repo_impl"
	"backend-github-trending/router"
	"github.com/labstack/echo"
)

func main() {

	sql := &db.Sql{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "092235",
		Dbname:   "golang",
	}
	sql.Connect()
	defer sql.Close() // thuc hien khi main chay xong moi chay no

	e := echo.New()
	structValidator := helper.NewStructValidator()
	structValidator.RegisterValidate()
	e.Validator = structValidator

	userHandler := handler.UserHandler{
		UserRepo: repo_impl.NewUserRepo(sql),
	}

	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
	}

	api.SetupRoutes()
	e.Logger.Fatal(e.Start(":3000"))
}
