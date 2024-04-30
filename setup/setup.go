package setup

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/THPTUHA/repeatword/config"
	"github.com/THPTUHA/repeatword/db"
)

//go:embed mysql/schema.sql
var mysqlSchema string

//go:embed mysql/functions.sql
var mysqlFunction string

//go:embed mysql/clean.sql
var mysqlClean string

func Setup() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg, err := config.Set(path.Join(dir, "config.yaml"))
	if err != nil {
		return err
	}

	cmd := exec.Command("mysql", "-u", cfg.DB.Mysql.Username, "-h", cfg.DB.Mysql.URI, "-P", fmt.Sprint(cfg.DB.Mysql.Port), "-p"+cfg.DB.Mysql.Password, "-e", fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", cfg.DB.Mysql.DatabaseName))

	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Error executing MySQL command: ", err)
	}

	sql, err := db.ConnectMysql()
	if err != nil {
		return err
	}
	schemas := strings.Split(mysqlSchema, ";")
	for _, sc := range schemas {
		if strings.TrimSpace(sc) == "" {
			continue
		}
		_, err = sql.ExecContext(context.Background(), sc)
		if err != nil {
			return err
		}
	}

	cleans := strings.Split(mysqlClean, ";")
	for _, cl := range cleans {
		if strings.TrimSpace(cl) == "" {
			continue
		}
		_, err = sql.ExecContext(context.Background(), cl)
		if err != nil {
			return err
		}
	}

	funtions := strings.Split(mysqlFunction, "-- Function")
	for _, fc := range funtions {
		if strings.TrimSpace(fc) == "" {
			continue
		}
		_, err = sql.ExecContext(context.Background(), fc)
		if err != nil {
			return err
		}
	}

	return nil
}
