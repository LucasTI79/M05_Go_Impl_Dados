package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	StorageDB *sql.DB
)

func initDatabase() {
	datasource := fmt.Sprintf(
		// user:pass@tcp(host:port)/dbname
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	var err error
	StorageDB, err = sql.Open("mysql", datasource)
	if err != nil {
		panic(err)
	}
	if err = StorageDB.Ping(); err != nil {
		panic(err)
	}

	log.Println("database Configured")
}
