package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

var MemeDB *sqlx.DB

type NullString struct {
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

func ConnectDB() {
	dbHost, dbPort, dbUsername, dbPassword := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS")
	database, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/meme", dbUsername, dbPassword, dbHost, dbPort))
	if err != nil {
		panic(err.Error())
		return
	}

	MemeDB = database
}
