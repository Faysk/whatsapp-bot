package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config é a estrutura que armazena configurações globais do sistema
type Config struct {
	DatabasePath string
	LogLevel     string
}

var AppConfig Config

// Load carrega as configurações do .env (se existir) e define os valores padrões
func Load() {
	// Carrega variáveis do .env, se o arquivo existir
	_ = godotenv.Load()

	AppConfig = Config{
		DatabasePath: getEnv("DB_PATH", "file:session.db?_pragma=foreign_keys(1)"),
		LogLevel:     getEnv("LOG_LEVEL", "INFO"),
	}

	log.Printf("📦 Configurações carregadas (DB_PATH=%s, LOG_LEVEL=%s)", AppConfig.DatabasePath, AppConfig.LogLevel)
}

// getEnv tenta buscar do sistema, se não achar retorna o padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
