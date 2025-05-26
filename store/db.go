package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/faysk/whatsapp-bot/config"
	_ "modernc.org/sqlite"
)

// InitDB cria e configura a conexão com o SQLite
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", config.AppConfig.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o banco de dados: %w", err)
	}

	// Evita problemas com bloqueios do SQLite em concorrência
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Testa a conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	log.Println("📁 Banco de dados conectado com sucesso.")
	return db, nil
}
