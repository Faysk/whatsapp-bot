package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/faysk/whatsapp-bot/config"
	"github.com/faysk/whatsapp-bot/store"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// InitWhatsAppClient inicializa o cliente do WhatsApp com sessão persistente via PostgreSQL
func InitWhatsAppClient(ctx context.Context) (*whatsmeow.Client, error) {
	logger := waLog.Stdout(config.AppConfig.BotName, config.AppConfig.LogLevel, true)

	driver := strings.ToLower(config.AppConfig.DatabaseDriver)
	dsn := config.AppConfig.DatabaseDSN

	if driver != "postgres" && driver != "postgresql" {
		return nil, fmt.Errorf("❌ Driver de banco de dados não suportado: %s", driver)
	}

	db, err := store.ConnectPostgres(dsn)
	if err != nil {
		return nil, fmt.Errorf("❌ Erro ao conectar ao PostgreSQL: %w", err)
	}

	container := sqlstore.NewWithDB(db, "postgres", logger)

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Println("⚠️ Nenhuma sessão ativa encontrada. Criando novo dispositivo...")
		deviceStore = container.NewDevice()
	}

	client := whatsmeow.NewClient(deviceStore, logger)
	log.Printf("✅ Cliente WhatsApp [%s] pronto para conectar.", config.AppConfig.BotName)
	return client, nil
}

// ConnectWithQR conecta com sessão existente ou exibe QR Code para pareamento manual
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
			return fmt.Errorf("❌ Erro durante escaneamento do QR")
		}
	}

	return nil
}
