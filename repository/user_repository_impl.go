package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iqbaltaufiq/latihan-restapi/helper"
	"github.com/iqbaltaufiq/latihan-restapi/model/domain"
)

type UserRepositoryImpl struct {
}

// create a constructor
// that will be called in main.go
func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

// insert user into table user.
// Save takes user object from service to be inserted into database
func (r *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	// lakukan query ke DB
	sql := "INSERT INTO user(name, occupation) VALUES (?,?)"
	result, err := tx.ExecContext(ctx, sql, user.Name, user.Occupation)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	user.Id = int(id)
	return user
}

// update user's name
func (r *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	sql := "UPDATE user SET name = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, sql, user.Name, user.Id)
	helper.PanicIfError(err)

	return user
}

// delete a user
func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, userId int) {
	sql := "DELETE FROM user WHERE id = ?"
	_, err := tx.ExecContext(ctx, sql, userId)
	helper.PanicIfError(err)
}

// get user by id
func (r *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, userId int) (domain.User, error) {
	sql := "SELECT * FROM user WHERE id = ?"
	rows, err := tx.QueryContext(ctx, sql, userId)
	helper.PanicIfError(err)

	user := domain.User{}

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Occupation)
		helper.PanicIfError(err)

		return user, nil
	} else {
		return user, errors.New("RepositoryError: User not found")
	}
}

// get all users
func (r *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.User {
	sql := "SELECT * FROM user"
	rows, err := tx.QueryContext(ctx, sql)
	helper.PanicIfError(err)

	users := []domain.User{}

	defer rows.Close()
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Occupation)
		helper.PanicIfError(err)

		users = append(users, user)
	}

	return users
}
