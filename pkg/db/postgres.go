package db

import (
	"database/sql"
	"fmt"
)

const (
	ProductsTable = "products"
)

type ConnectionInfo struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConnection(cfg ConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.SSLMode, cfg.Password))
	if err != nil{
		return nil, err
	}

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return db, nil
}
