package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/faysk/whatsapp-bot/config"
	"github.com/faysk/whatsapp-bot/store"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

const fallbackDSN = "file:session.db?_pragma=foreign_keys(1)&_journal_mode=WAL&_busy_timeout=5000"

// InitWhatsAppClient inicializa o cliente do WhatsApp com sessão persistente (PostgreSQL ou SQLite)
func InitWhatsAppClient(ctx context.Context) (*whatsmeow.Client, error) {
	logger := waLog.Stdout(config.AppConfig.BotName, config.AppConfig.LogLevel, true)

	driver := strings.ToLower(config.AppConfig.DatabaseDriver)
	dsn := config.AppConfig.DatabasePath

	var container *sqlstore.Container
	var err error

	switch driver {
	case "postgres", "postgresql":
		var db *sql.DB
		db, err = store.ConnectPostgres(dsn)
		if err != nil {
			return nil, fmt.Errorf("❌ Erro ao conectar com PostgreSQL: %w", err)
		}
		container = sqlstore.NewWithDB(db, "postgres", logger)
	default:
		dsn = sanitizeDSN(dsn)
		var db *sql.DB
		db, err = sql.Open("sqlite", dsn)
		if err != nil {
			return nil, fmt.Errorf("❌ Erro ao abrir banco SQLite: %w", err)
		}
		if err := db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("❌ Falha ao conectar ao SQLite: %w", err)
		}
		container = sqlstore.NewWithDB(db, "sqlite", logger)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Println("⚠️ Nenhuma sessão ativa encontrada. Criando novo dispositivo...")
		deviceStore = container.NewDevice()
	}

	client := whatsmeow.NewClient(deviceStore, logger)
	log.Printf("✅ Cliente WhatsApp [%s] pronto para conectar.", config.AppConfig.BotName)
	return client, nil
}

// ConnectWithQR realiza a conexão ou mostra o QR Code caso não haja sessão
func ConnectWithQR(ctx context.Context, client *whatsmeow.Client) error {
	if client.Store.ID != nil {
		log.Println("🔗 Reconectando com sessão existente...")
		return client.Connect()
	}

	log.Println("📱 Iniciando pareamento via QR Code...")

	qrChan, _ := client.GetQRChannel(ctx)
	if err := client.Connect(); err != nil {
		return fmt.Errorf("❌ Falha ao conectar: %w", err)
	}

	for evt := range qrChan {
		switch evt.Event {
		case "code":
			fmt.Println("📷 Escaneie o QR abaixo para parear:")
			qr := qrcodeTerminal.New()
			qr.Get(evt.Code).Print()
		case "success":
			log.Println("✅ QR Code escaneado com sucesso!")
			return nil
		case "timeout":
			return fmt.Errorf("⏳ Tempo esgotado para escanear o QR")
		case "error":
			return fmt.Errorf("❌ Erro ao escanear o QR")
		}
	}

	return nil
}

// sanitizeDSN garante que o DSN tenha os parâmetros essenciais de concorrência
func sanitizeDSN(dsn string) string {
	if strings.Contains(dsn, "_journal_mode") && strings.Contains(dsn, "_busy_timeout") {
		return dsn
	}
	log.Println("⚠️ DSN incompleto no .env — aplicando padrão seguro para concorrência (WAL + timeout)")
	return fallbackDSN
}
