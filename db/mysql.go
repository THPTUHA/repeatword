package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectMysql() (*sql.DB, error) {
	var err error
	if db != nil {
		return db, nil
	}
	username := "nghia"
	password := "root"
	hostname := "127.0.0.1"
	dbname := "repeatword"

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	return db, err
}
