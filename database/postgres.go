package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type ServiceHandler struct {
	DB *sql.DB
}

func Connect(conStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
