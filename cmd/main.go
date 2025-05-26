package main

import (
	"context"
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
)

func main() {
	log.Println("🚀 Iniciando o bot WhatsApp...")

	// 📝 Inicializa logger (timestamp + cores)
	utils.SetupLogger()

	// 🔧 Carrega configurações do .env e popula AppConfig
	config.Load()

	// 🔐 Carrega números autorizados dinâmicos (JSON)
	dynamic := store.LoadAuthorizedNumbers()
	config.AddDynamicAuthorizedNumbers(dynamic)

	// 🌐 Cria contexto principal para uso compartilhado
	ctx := context.Background()

	// 📲 Inicializa cliente WhatsApp com persistência no banco
	client, err := services.InitWhatsAppClient(ctx)
	if err != nil {
		log.Fatalf("❌ Falha ao inicializar o cliente WhatsApp: %v", err)
	}

	// 🔗 Conecta com sessão existente ou via QR Code
	if err := services.ConnectWithQR(ctx, client); err != nil {
		log.Fatalf("❌ Erro na conexão com WhatsApp: %v", err)
	}

	log.Println("✅ Bot conectado com sucesso. Aguardando mensagens...")

	// 🗞️ Agendador de notícias cripto (diário às 10h)
	scheduler.StartDailyNews(ctx, client, config.AppConfig.AuthorizedNumbers)

	// 📩 Listener de eventos WhatsApp (mensagens recebidas)
	events.Listen(ctx, client)

	// ⛔ Finalização segura com Ctrl+C ou SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("📴 Encerrando conexão com o WhatsApp...")
	client.Disconnect()
}
