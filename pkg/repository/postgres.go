package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	advertsTable = "adverts"
)

type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	fmt.Println(cfg)
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.UserName, cfg.Password, cfg.Host, cfg.DBName, cfg.SSLMode,
	))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}