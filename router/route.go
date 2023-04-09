package router

import (
	"github.com/iqbaltaufiq/latihan-restapi/controller"
	"github.com/iqbaltaufiq/latihan-restapi/exception"
	"github.com/julienschmidt/httprouter"
)

// write down all of your routes here
func NewRouter(controller controller.UserController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/users", controller.FindAll)
	router.GET("/api/users/:userId", controller.FindById)
	router.POST("/api/users", controller.Create)
	router.PUT("/api/users/:userId", controller.Update)
	router.DELETE("/api/users/:userId", controller.Delete)

	router.PanicHandler = exception.PanicHandler
	return router
}
