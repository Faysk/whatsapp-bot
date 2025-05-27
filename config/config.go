package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

//
// ========== 🔧 Estrutura =========
//

// Config representa todas as configurações carregadas do .env
type Config struct {
	DatabaseDriver     string
	DatabasePath       string
	LogLevel           string
	Port               string
	OpenAIKey          string
	OpenAIModel        string
	EnableChatGPT      bool
	BotName            string
	Language           string
	MaxTokens          int
	Temperature        float64
	RestrictToGroup    bool
	FixedAuthorizedEnv []string
	AuthorizedNumbers  []string
}

// AppConfig é a instância global acessada pelo projeto
var AppConfig Config

//
// ========== 🚀 Carregamento =========
//

// Load carrega as variáveis do .env e preenche o AppConfig
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Erro ao carregar .env: %v", err)
	}

	fixedNumbers := parseCSVEnv("AUTHORIZED_NUMBERS")

	AppConfig = Config{
		DatabaseDriver:     getEnv("DB_DRIVER", "sqlite"),
		DatabasePath:       getEnv("DB_PATH", "file:session.db?_pragma=foreign_keys(1)"),
		LogLevel:           getEnv("LOG_LEVEL", "INFO"),
		Port:               getEnv("PORT", "8080"),
		OpenAIKey:          getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:        getEnv("OPENAI_MODEL", "gpt-4o"),
		EnableChatGPT:      getBool("ENABLE_CHATGPT", true),
		BotName:            getEnv("BOT_NAME", "FayskBot"),
		Language:           getEnv("LANG", "pt-BR"),
		MaxTokens:          getInt("MAX_TOKENS", 400),
		Temperature:        getFloat("TEMPERATURE", 0.7),
		RestrictToGroup:    getBool("RESTRICT_TO_GROUP", false),
		FixedAuthorizedEnv: fixedNumbers,
		AuthorizedNumbers:  append([]string{}, fixedNumbers...),
	}

	printLoadedConfig()
}

// printLoadedConfig exibe todas as configurações carregadas
func printLoadedConfig() {
	log.Println("📦 Configurações carregadas:")
	log.Printf("  ├─ DB_DRIVER:          %s", AppConfig.DatabaseDriver)
	log.Printf("  ├─ DB_PATH:            %s", AppConfig.DatabasePath)
	log.Printf("  ├─ LOG_LEVEL:          %s", AppConfig.LogLevel)
	log.Printf("  ├─ PORT:               %s", AppConfig.Port)
	log.Printf("  ├─ BOT_NAME:           %s", AppConfig.BotName)
	log.Printf("  ├─ LANG:               %s", AppConfig.Language)
	log.Printf("  ├─ OPENAI_MODEL:       %s", AppConfig.OpenAIModel)
	log.Printf("  ├─ MAX_TOKENS:         %d", AppConfig.MaxTokens)
	log.Printf("  ├─ TEMPERATURE:        %.2f", AppConfig.Temperature)
	log.Printf("  ├─ RESTRICT_TO_GROUP:  %v", AppConfig.RestrictToGroup)
	log.Printf("  ├─ FIXED NUMBERS:      %v", AppConfig.FixedAuthorizedEnv)

	if AppConfig.OpenAIKey != "" && AppConfig.EnableChatGPT {
		log.Println("  └─ IA: ✅ habilitada (ChatGPT ativo)")
	} else if AppConfig.OpenAIKey == "" {
		log.Println("  └─ IA: ⚠️ desabilitada - OPENAI_API_KEY ausente")
	} else {
		log.Println("  └─ IA: ⚠️ desabilitada via ENABLE_CHATGPT=false")
	}
}

//
// ========== ➕ Autorização dinâmica =========
//

// AddDynamicAuthorizedNumbers adiciona números dinâmicos à lista final, sem duplicar os fixos
func AddDynamicAuthorizedNumbers(dynamic []string) {
	for _, n := range dynamic {
		if !contains(AppConfig.FixedAuthorizedEnv, n) && !contains(AppConfig.AuthorizedNumbers, n) {
			AppConfig.AuthorizedNumbers = append(AppConfig.AuthorizedNumbers, n)
		}
	}
}

//
// ========== 🧰 Utilitários =========
//

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}

func getFloat(key string, defaultValue float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

func getBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultValue
}

func parseCSVEnv(key string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}
