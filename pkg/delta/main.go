package main

import (
	"fmt"

	"github.com/go-Echelon/go-Echelon/pkg/core/database"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	dbName := "mydb"
	db, err := database.Connect(mongoURI, dbName)
	if err != nil {
		fmt.Println("ok now there is no error")
	}

	fmt.Println("got the db", db)
}
