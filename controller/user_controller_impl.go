package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/iqbaltaufiq/latihan-restapi/helper"
	"github.com/iqbaltaufiq/latihan-restapi/model/web"
	"github.com/iqbaltaufiq/latihan-restapi/service"
	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService service.UserService
}

// create a constructor
// that will be called in main.go
func NewUserController(UserService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: UserService,
	}
}

func (c *UserControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	payload := web.UserCreatePayload{}

	// decode request
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&payload)

	// send it to service
	serviceResponse := c.UserService.Create(request.Context(), payload)

	// send the returned value into encoder
	response := web.HttpResponse{
		Code:   200,
		Status: "OK",
		Data:   serviceResponse,
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	helper.PanicIfError(err)
}

func (c *UserControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// prepare an empty container
	payload := web.UserUpdatePayload{}

	// decode the request stream
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&payload)

	// userId in param is a string
	// convert it to int first
	userId, err := strconv.Atoi(params.ByName("userId"))
	helper.PanicIfError(err)

	payload.Id = userId

	serviceResponse := c.UserService.Update(request.Context(), payload)

	response := web.HttpResponse{
		Code:   200,
		Status: "OK",
		Data:   serviceResponse,
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	encoder.Encode(response)
}

func (c *UserControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, err := strconv.Atoi(params.ByName("userId"))
	helper.PanicIfError(err)

	c.UserService.Delete(request.Context(), userId)

	response := web.HttpResponse{
		Code:   200,
		Status: "OK",
		Data:   "Deleted successfully",
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(response)
	helper.PanicIfError(err)
}

func (c *UserControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, err := strconv.Atoi(params.ByName("userId"))
	helper.PanicIfError(err)

	user := c.UserService.FindById(request.Context(), userId)
	response := web.HttpResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(response)
	helper.PanicIfError(err)
}

func (c *UserControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	users := c.UserService.FindAll(request.Context())

	response := web.HttpResponse{
		Code:   200,
		Status: "OK",
		Data:   users,
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	helper.PanicIfError(err)
}
