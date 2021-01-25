package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var MemeDB *sql.DB

func ConnectDB() {
	dbHost, dbPort, dbUsername, dbPassword := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS")
	database, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/meme", dbUsername, dbPassword, dbHost, dbPort))
	if err != nil {
		panic(err.Error())
		return
	}

	MemeDB = database
}
