package database

import (
	"database/sql"
	"fmt"

	"github.com/hashicorp/go-hclog"
	_ "github.com/lib/pq"
)

type DBSettings struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func InitDB(settings *DBSettings, logger hclog.Logger) *sql.DB {
	connStr := fmt.Sprint(
		"postgresql://",
		settings.User,
		":",
		settings.Password,
		"@",
		settings.Host,
		":",
		settings.Port,
		"/",
		settings.DBName,
		"?sslmode=disable",
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Error starting db", "error", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Error connecting to db", "error", err)
	}

	if err == nil {
		logger.Info("Database connection initialized...")
	}

	defer db.Close()

	return db
}
