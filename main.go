// @author Iqbal Taufiqurrahman
// created in April 2023
// This is a repository created after completing
// Programmer Zaman Now's Golang RESTful API course
// in Codepolitan.
// This repository is for educational purposes and free of use.

package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iqbaltaufiq/latihan-restapi/app"
	"github.com/iqbaltaufiq/latihan-restapi/controller"
	"github.com/iqbaltaufiq/latihan-restapi/middleware"
	"github.com/iqbaltaufiq/latihan-restapi/repository"
	"github.com/iqbaltaufiq/latihan-restapi/router"
	"github.com/iqbaltaufiq/latihan-restapi/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate)
	userController := controller.NewUserController(userService)

	httpRouter := router.NewRouter(userController)

	// apply auth middleware in all routes
	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(httpRouter),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
