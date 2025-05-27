package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectPostgres cria uma conexão com PostgreSQL
func ConnectPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no PostgreSQL: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}
