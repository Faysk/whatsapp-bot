package store

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/faysk/whatsapp-bot/config"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// InitClient inicializa o client WhatsApp com persistência no SQLite
func InitClient(ctx context.Context) (*whatsmeow.Client, error) {
	db, err := InitDB()
	if err != nil {
		return nil, err
	}

	logger := waLog.Stdout("WhatsAppBot", config.AppConfig.LogLevel, true)
	container := sqlstore.NewWithDB(db, "sqlite", logger)

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Println("⚠️ Nenhuma sessão ativa encontrada. Iniciando pareamento via QR Code...")
		deviceStore = container.NewDevice()
	}

	client := whatsmeow.NewClient(deviceStore, logger)

	if err := connectClient(ctx, client); err != nil {
		return nil, fmt.Errorf("erro ao conectar cliente WhatsApp: %w", err)
	}

	return client, nil
}

// connectClient realiza conexão com o WhatsApp e mostra QR Code se necessário
func connectClient(ctx context.Context, client *whatsmeow.Client) error {
	if client.Store.ID != nil {
		return client.Connect()
	}

	qrChan, _ := client.GetQRChannel(ctx)
	if err := client.Connect(); err != nil {
		return err
	}

	for evt := range qrChan {
		switch evt.Event {
		case "code":
			fmt.Println("📷 Escaneie o QR abaixo para parear:")
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		case "success":
			log.Println("✅ QR Code escaneado com sucesso!")
			return nil
		case "timeout":
			return fmt.Errorf("⏳ Tempo esgotado para escanear o QR")
		case "error":
			return fmt.Errorf("❌ Erro ao escanear o QR Code")
		}
	}
	return nil
}
