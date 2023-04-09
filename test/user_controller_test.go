package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iqbaltaufiq/latihan-restapi/controller"
	"github.com/iqbaltaufiq/latihan-restapi/helper"
	"github.com/iqbaltaufiq/latihan-restapi/middleware"
	"github.com/iqbaltaufiq/latihan-restapi/model/domain"
	"github.com/iqbaltaufiq/latihan-restapi/repository"
	"github.com/iqbaltaufiq/latihan-restapi/router"
	"github.com/iqbaltaufiq/latihan-restapi/service"
	"github.com/stretchr/testify/assert"
)

// This is an integration testing to test all of the functions
// in the User controller.
// To make sure that all functions return the desired response.

// WARNING:
// make sure to change the database to database for testing
func setupDBTest() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/latihan_go_restapi_test")
	helper.PanicIfError(err)

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate)
	userController := controller.NewUserController(userService)

	router := router.NewRouter(userController)

	return middleware.NewAuthMiddleware(router)
}

// truncate the table whenever you run a test
// so we always run a test with an empty db
func truncateDB(db *sql.DB) {
	db.Exec("TRUNCATE user")
}

func TestCreateUserSuccess(t *testing.T) {
	// intialize db and router
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	// create a payload as request body
	payload := strings.NewReader(`{"name": "John", "occupation": "student"}`)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/users", payload)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	var responseBody map[string]interface{}
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "John", responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateUserFailed(t *testing.T) {
	// make a db connection
	// make sure to use database for testing
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	payload := strings.NewReader(`{"name": "barusu, "occupation": "student"}`)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/users", payload)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	var responseBody map[string]interface{}
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "Bad Request", responseBody["status"])
}

func TestUpdateUserSuccess(t *testing.T) {
	// make a database connection
	// make sure to use database for testing
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	// make a db transaction.
	// insert an entry into db
	// that will be updated afterwards.
	tx, _ := db.Begin()
	user := repository.NewUserRepository().Save(context.Background(), tx, domain.User{
		Name:       "John Doe",
		Occupation: "student",
	})
	tx.Commit()

	fmt.Println(user)

	// create a payload to be sent into request
	payloadJSON, _ := json.Marshal(domain.User{
		Name: "Jack",
	})

	payload := strings.NewReader(string(payloadJSON))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/users/"+strconv.Itoa(user.Id), payload)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	fmt.Print(responseBody)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, user.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
}

func TestUpdateUserFailed(t *testing.T) {
	// make a database connection
	// make sure to use database for testing
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	// make a db transaction.
	// insert an entry into db
	// that will be updated afterwards.
	tx, _ := db.Begin()
	user := repository.NewUserRepository().Save(context.Background(), tx, domain.User{
		Name:       "John Doe",
		Occupation: "student",
	})
	tx.Commit()

	fmt.Println(user)

	// create a payload to be sent into request
	payloadJSON, _ := json.Marshal(domain.User{
		Name: "",
	})

	payload := strings.NewReader(string(payloadJSON))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/users/"+strconv.Itoa(user.Id), payload)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	fmt.Print(responseBody)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "Bad Request", responseBody["status"])
}

func TestFindUserSuccess(t *testing.T) {
	// make a database connection
	// make sure to use database for testing purposes only
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	tx, _ := db.Begin()
	user := repository.NewUserRepository().Save(context.Background(), tx, domain.User{
		Name:       "John",
		Occupation: "student",
	})
	tx.Commit()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/users/"+strconv.Itoa(user.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "John", responseBody["data"].(map[string]interface{})["name"])
}

func TestFindUserFailed(t *testing.T) {
	// make a database connection
	// make sure to use database for testing purposes only
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/users/100", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, response.StatusCode)
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
}

func TestDeleteUserSuccess(t *testing.T) {
	// make a database connection
	// make sure to use database for testing purposes only
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	// make a db transaction
	tx, _ := db.Begin()
	user := repository.NewUserRepository().Save(context.Background(), tx, domain.User{
		Name:       "John",
		Occupation: "student",
	})
	tx.Commit()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/users/"+strconv.Itoa(user.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteUserFailed(t *testing.T) {
	// make a database connection
	// make sure to use database for testing purposes only
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/users/100", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, response.StatusCode)
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "Not Found", responseBody["status"])
}

func TestFindUsers(t *testing.T) {
	// make a connection to db
	// make sure to use database for testing purposes only
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	// make a db transaction
	// to insert few users
	tx, _ := db.Begin()
	user1 := repository.NewUserRepository().Save(context.Background(), tx, domain.User{
		Name:       "John",
		Occupation: "student",
	})
	user2 := repository.NewUserRepository().Save(context.Background(), tx, domain.User{
		Name:       "Anne",
		Occupation: "lecturer",
	})
	tx.Commit()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/users", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SECRET")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, user1.Name, responseBody["data"].([]interface{})[0].(map[string]interface{})["name"])
	assert.Equal(t, user2.Name, responseBody["data"].([]interface{})[1].(map[string]interface{})["name"])
}

func TestUnauthorized(t *testing.T) {
	// make a db connection
	// make sure to use database for testing purposes only
	db := setupDBTest()
	truncateDB(db)

	router := setupRouter(db)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/users", nil)
	request.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 401, response.StatusCode)
	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "Unauthorized", responseBody["status"])
}
