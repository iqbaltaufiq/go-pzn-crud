package router

import (
	"github.com/iqbaltaufiq/latihan-restapi/controller"
	"github.com/julienschmidt/httprouter"
)

// write down all of your routes here
func NewRouter(controller controller.UserController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/categories", controller.FindAll)
	router.GET("/api/categories/:userId", controller.FindById)
	router.POST("/api/categories", controller.Create)
	router.PUT("/api/categories/:userId", controller.Update)
	router.DELETE("/api/categories/:userId", controller.Delete)

	return router
}
