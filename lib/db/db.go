package db

import (
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Connect(dsn string) (*sqlx.DB, error) {

	start := time.Now()
	log.Println("Connecting to db")

	init := "SET SESSION SQL_MODE='STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ONLY_FULL_GROUP_BY'"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// sqlx needs to know which columns to map to a struct
	// it does case sensitive matching, and we use lowerCamelCase for the db
	// but UpperCamelCase for structs (so golang will export them)
	sqlx.NameMapper = upperCamelToLowerCamel

	_, err = db.Exec(init)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to db in %s", time.Since(start))

	return db, nil

}

func upperCamelToLowerCamel(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
