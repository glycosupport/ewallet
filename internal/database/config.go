package database

import (
	"os"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		username: os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		database: os.Getenv("DB_DATABASE"),
	}
}
