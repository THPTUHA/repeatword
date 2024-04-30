package db

import (
	"database/sql"
	"fmt"

	"github.com/THPTUHA/repeatword/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectMysql() (*sql.DB, error) {
	var err error
	if db != nil {
		return db, nil
	}
	var username, password, hostname, dbname string
	config, _ := config.Get()
	if config == nil {
		username = "nghia"
		password = "root"
		hostname = "127.0.0.1"
		dbname = "repeatword"
	} else {
		username = config.DB.Mysql.Username
		password = config.DB.Mysql.Password
		hostname = config.DB.Mysql.URI
		dbname = config.DB.Mysql.DatabaseName
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	return db, err
}
