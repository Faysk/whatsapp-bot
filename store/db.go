package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// ConnectPostgres cria uma conexão com PostgreSQL com verificação ativa
func ConnectPostgres(dsn string) (*sql.DB, error) {
	start := time.Now()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("❌ Falha ao abrir conexão com PostgreSQL: %w", err)
	}

	// Testa a conexão antes de retornar
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("❌ Conexão com PostgreSQL falhou: %w", err)
	}

	// Configurações recomendadas para controle de conexões
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Printf("✅ Conexão com PostgreSQL estabelecida em %s", time.Since(start))
	return db, nil
}
