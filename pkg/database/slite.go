package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDB(datasourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", datasourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
