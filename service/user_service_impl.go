package service

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/iqbaltaufiq/latihan-restapi/helper"
	"github.com/iqbaltaufiq/latihan-restapi/model/domain"
	"github.com/iqbaltaufiq/latihan-restapi/model/web"
	"github.com/iqbaltaufiq/latihan-restapi/repository"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

// create a constructor
// that will be called in main.go
func NewUserService(UserRepository repository.UserRepository, DB *sql.DB, Validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: UserRepository,
		DB:             DB,
		Validate:       Validate,
	}
}

func (s *UserServiceImpl) Create(ctx context.Context, request web.UserCreatePayload) web.UserResponse {
	// do validation for the payload struct
	err := s.Validate.Struct(request)
	helper.PanicIfError(err)

	// start a db transaction
	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// payload mapping before being sent to repository
	payload := domain.User{
		Name:       request.Name,
		Occupation: request.Occupation,
	}

	// send payload to repository
	// to be inserted into DB
	user := s.UserRepository.Save(ctx, tx, payload)

	return web.UserResponse{
		Id:         user.Id,
		Name:       user.Name,
		Occupation: user.Occupation,
	}
}

func (s *UserServiceImpl) Update(ctx context.Context, request web.UserUpdatePayload) web.UserResponse {
	// do validation for the payload struct
	err := s.Validate.Struct(request)
	helper.PanicIfError(err)

	// make a db transaction
	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// find the user in DB
	userInDB, err := s.UserRepository.FindById(ctx, tx, request.Id)
	helper.PanicIfError(err)

	userInDB.Name = request.Name

	user := s.UserRepository.Update(ctx, tx, userInDB)
	return web.UserResponse{
		Id:         user.Id,
		Name:       user.Name,
		Occupation: user.Occupation,
	}
}

func (s *UserServiceImpl) Delete(ctx context.Context, userId int) {
	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(ctx, tx, userId)
	helper.PanicIfError(err)

	s.UserRepository.Delete(ctx, tx, user.Id)
}

func (s *UserServiceImpl) FindById(ctx context.Context, userId int) web.UserResponse {
	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(ctx, tx, userId)
	helper.PanicIfError(err)

	return web.UserResponse{
		Id:         user.Id,
		Name:       user.Name,
		Occupation: user.Occupation,
	}
}

func (s *UserServiceImpl) FindAll(ctx context.Context) []web.UserResponse {
	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	users := s.UserRepository.FindAll(ctx, tx)

	var responses []web.UserResponse
	for _, user := range users {
		responses = append(responses, web.UserResponse{
			Id:         user.Id,
			Name:       user.Name,
			Occupation: user.Occupation,
		})
	}

	return responses
}
