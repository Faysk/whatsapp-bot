package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/faysk/whatsapp-bot/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "modernc.org/sqlite"
)

// InitWhatsAppClient inicializa e retorna o cliente WhatsApp com sessão persistente
func InitWhatsAppClient(ctx context.Context) (*whatsmeow.Client, error) {
	logger := waLog.Stdout(config.AppConfig.BotName, config.AppConfig.LogLevel, true)

	db, err := sql.Open("sqlite", config.AppConfig.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("❌ Erro ao abrir banco SQLite: %w", err)
	}

	if pingErr := db.PingContext(ctx); pingErr != nil {
		return nil, fmt.Errorf("❌ Falha ao conectar ao banco SQLite: %w", pingErr)
	}

	container := sqlstore.NewWithDB(db, "sqlite", logger)

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Println("⚠️ Nenhuma sessão ativa encontrada. Criando novo dispositivo...")
		deviceStore = container.NewDevice()
	}

	client := whatsmeow.NewClient(deviceStore, logger)
	log.Printf("✅ Cliente WhatsApp [%s] pronto para conectar.", config.AppConfig.BotName)
	return client, nil
}

// ConnectWithQR conecta o cliente ao WhatsApp e lida com pareamento via QR Code
func ConnectWithQR(ctx context.Context, client *whatsmeow.Client) error {
	if client.Store.ID != nil {
		log.Println("🔗 Reconectando com sessão existente...")
		return client.Connect()
	}

	log.Println("📱 Iniciando pareamento via QR Code...")

	qrChan, _ := client.GetQRChannel(ctx)
	if err := client.Connect(); err != nil {
		return fmt.Errorf("falha ao conectar: %w", err)
	}

	for evt := range qrChan {
		switch evt.Event {
		case "code":
			fmt.Println("📷 Escaneie o QR abaixo para parear:")
			fmt.Println(evt.Code)
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
