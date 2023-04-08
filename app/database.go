package app

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iqbaltaufiq/latihan-restapi/helper"
)

// create a database connection and set pooling
// dsn format : user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/latihan_go_restapi")
	helper.PanicIfError(err)

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
