package service

import (
	"context"

	"github.com/iqbaltaufiq/latihan-restapi/model/web"
)

type UserService interface {
	Create(ctx context.Context, request web.UserCreatePayload) web.UserResponse
	Update(ctx context.Context, request web.UserUpdatePayload) web.UserResponse
	Delete(ctx context.Context, userId int)
	FindById(ctx context.Context, userId int) web.UserResponse
	FindAll(ctx context.Context) []web.UserResponse
}
