package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/faysk/whatsapp-bot/config"
	"github.com/faysk/whatsapp-bot/events"
	"github.com/faysk/whatsapp-bot/scheduler"
	"github.com/faysk/whatsapp-bot/services"
	"github.com/faysk/whatsapp-bot/store"
	"github.com/faysk/whatsapp-bot/utils"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("🔥 Pânico recuperado: %v", r)
		}
	}()

	log.Println("🚀 Iniciando o bot WhatsApp...")
	utils.SetupLogger()
	ctx := context.Background()

	client, err := startBot(ctx)
	if err != nil {
		log.Fatalf("❌ Erro crítico: %v", err)
	}

	log.Println("✅ Bot conectado com sucesso. Aguardando mensagens...")

	scheduler.StartDailyNews(ctx, client, config.AppConfig.AuthorizedNumbers)

	services.MonitorCryptos(func(msg string) {
		for _, number := range config.AppConfig.AuthorizedNumbers {
			jid := types.NewJID(number, "s.whatsapp.net")
			services.SendReply(ctx, client, jid, msg)
		}
	})

	events.Listen(ctx, client)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("📴 Encerrando conexão com o WhatsApp...")
	if client.IsConnected() {
		client.Disconnect()
	}
}

func startBot(ctx context.Context) (*whatsmeow.Client, error) {
	config.Load()

	dynamic := store.LoadAuthorizedNumbers()
	config.AddDynamicAuthorizedNumbers(dynamic)

	client, err := services.InitWhatsAppClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar cliente WhatsApp: %w", err)
	}

	if err := services.ConnectWithQR(ctx, client); err != nil {
		return nil, fmt.Errorf("erro ao conectar com QR: %w", err)
	}

	return client, nil
}
