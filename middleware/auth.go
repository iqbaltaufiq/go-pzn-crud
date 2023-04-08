package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/iqbaltaufiq/latihan-restapi/helper"
	"github.com/iqbaltaufiq/latihan-restapi/model/web"
)

type AuthMiddleware struct {
	Handler http.Handler
}

// make a constructor
// that will be called in main.go
func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

// make authentication middleware that checks for "X-API-KEY"
// this middleware will be placed in ALL routes
// the hardcoded API key is "SECRET"
func (m *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("X-API-KEY") != "SECRET" {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		response := web.HttpResponse{
			Code:   301,
			Status: "Unauthorized",
		}

		encoder := json.NewEncoder(writer)
		err := encoder.Encode(response)
		helper.PanicIfError(err)
	} else {
		m.Handler.ServeHTTP(writer, request)
	}
}
