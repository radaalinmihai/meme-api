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

type NullString sql.NullString

func (x * NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.String)
}

func (x *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		x.Valid = true
		x.String = *s
	} else {
		x.Valid = false
	}
	return nil
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
