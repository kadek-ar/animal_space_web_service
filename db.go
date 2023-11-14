package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func DatabaseConnection() {
	fmt.Println("connecting database")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")

	if err != nil {
		panic((err.Error()))
	}

	defer db.Close()

	insert, err := db.Query("INSERT INTO test VALUES ( 2, 'test name', 'test artis', 45.45 )")

	if err != nil {
		panic((err.Error()))
	}

	defer insert.Close()
}
