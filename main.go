package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/milhamap/simple-note/helper"
)

func main() {
	server := InitializedServer()

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
