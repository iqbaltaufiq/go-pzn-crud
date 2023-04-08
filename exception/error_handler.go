package exception

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iqbaltaufiq/latihan-restapi/helper"
	"github.com/iqbaltaufiq/latihan-restapi/model/web"
)

// Handler that will process panic.
// This will be called in router.PanicHandler.
func PanicHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	fmt.Println("=== Masuk panic handler ===")
	if notFoundError(writer, request, err) {
		return
	}

	if validationError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}

func validationError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)

	if !ok {
		return false
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)

	response := web.HttpResponse{
		Code:   http.StatusBadRequest,
		Status: "Bad Request",
		Data:   exception.Error(),
	}

	encoder := json.NewEncoder(writer)
	encodeErr := encoder.Encode(response)
	helper.PanicIfError(encodeErr)

	return true
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	fmt.Println("Masuk error handler")
	exception, ok := err.(NotFoundError)

	if !ok {
		return false
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNotFound)

	response := web.HttpResponse{
		Code:   http.StatusNotFound,
		Status: "Not Found",
		Data:   exception.Error,
	}

	encoder := json.NewEncoder(writer)
	encodeErr := encoder.Encode(response)
	helper.PanicIfError(encodeErr)

	return true
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	response := web.HttpResponse{
		Code:   http.StatusInternalServerError,
		Status: "Internal Server Error",
		Data:   err,
	}

	encoder := json.NewEncoder(writer)
	encodeErr := encoder.Encode(response)
	helper.PanicIfError(encodeErr)
}
