package app

import (
	"database/sql"
	"github.com/milhamap/simple-note/helper"
	"time"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/simple_note")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
	//migrate -database "mysql://root@tcp(localhost:3306)/simple_note" -path db/migrations up
}
