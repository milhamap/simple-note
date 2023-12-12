package main

import (
	"belajar-golang-restful-api/helper"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	server := InitializedServer()

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
