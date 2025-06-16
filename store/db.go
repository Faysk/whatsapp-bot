package store

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// ConnectPostgres cria uma conexão robusta com PostgreSQL, incluindo validação e tuning de pool
func ConnectPostgres(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("❌ DSN do PostgreSQL está vazio")
	}

	// Validação de prefixo
	if !strings.HasPrefix(dsn, "postgres://") && !strings.HasPrefix(dsn, "postgresql://") {
		return nil, fmt.Errorf("❌ DSN inválido, deve começar com postgres:// ou postgresql://")
	}

	parsed, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("❌ Erro ao parsear DSN: %w", err)
	}

	start := time.Now()
	log.Printf("🔌 Conectando ao banco PostgreSQL em %s (DB: %s)", parsed.Host, strings.TrimPrefix(parsed.Path, "/"))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("❌ Falha ao abrir conexão com PostgreSQL: %w", err)
	}

	// Pool tuning: permite override via ENV (útil em prod)
	db.SetMaxOpenConns(getEnvInt("PG_MAX_OPEN_CONNS", 10))
	db.SetMaxIdleConns(getEnvInt("PG_MAX_IDLE_CONNS", 5))
	db.SetConnMaxLifetime(30 * time.Minute)

	// Tentativas com backoff leve para containers lentos
	for i := 1; i <= 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("⏳ Tentativa %d de conexão com PostgreSQL falhou: %v", i, err)
		time.Sleep(time.Duration(i) * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("❌ Conexão com PostgreSQL falhou após tentativas: %w", err)
	}

	log.Printf("✅ Conexão com PostgreSQL estabelecida em %s", time.Since(start))
	return db, nil
}

// getEnvInt lê variável de ambiente como inteiro ou retorna fallback
func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	var i int
	if _, err := fmt.Sscanf(val, "%d", &i); err != nil {
		return fallback
	}
	return i
}
