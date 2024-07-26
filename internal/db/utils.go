package db

import (
	"database/sql"
	"log"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

func InitDB(dbName string, address string, user string, password string) {
	cfg := getConfig(dbName, address, user, password)

	var openErr error

	DB, openErr = sql.Open("mysql", cfg.FormatDSN())

	if openErr != nil {
		log.Fatalf("Error Initializing Database: %v", openErr)
	}

	slog.Info("Database Initialized")
}

func getConfig(dbName string, address string, user string, password string) mysql.Config {
	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   address,
		DBName: dbName,
	}

	return cfg
}
