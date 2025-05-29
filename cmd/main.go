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
	"go.mau.fi/whatsmeow/types"
)

func main() {
	log.Println("🚀 Iniciando o bot WhatsApp...")

	// 📝 Inicializa logger (timestamp + cores)
	utils.SetupLogger()

	// 🔧 Carrega configurações do .env e popula AppConfig
	config.Load()

	// 🔐 Carrega números autorizados dinâmicos do arquivo JSON
	dynamic := store.LoadAuthorizedNumbers()
	config.AddDynamicAuthorizedNumbers(dynamic)

	// 🌐 Cria contexto principal
	ctx := context.Background()

	// 📲 Inicializa cliente WhatsApp
	client, err := services.InitWhatsAppClient(ctx)
	if err != nil {
		log.Fatalf("❌ Falha ao inicializar o cliente WhatsApp: %v", err)
	}

	// 🔗 Conecta usando sessão persistida ou QR
	if err := services.ConnectWithQR(ctx, client); err != nil {
		log.Fatalf("❌ Erro ao conectar com WhatsApp: %v", err)
	}

	log.Println("✅ Bot conectado com sucesso. Aguardando mensagens...")

	// 🗞️ Agendador de notícias diárias sobre criptomoedas
	scheduler.StartDailyNews(ctx, client, config.AppConfig.AuthorizedNumbers)

	// 🚨 Inicia o monitor de recordes de criptoativos (ATH)
	services.MonitorCryptos(func(msg string) {
		for _, number := range config.AppConfig.AuthorizedNumbers {
			jid := types.NewJID(number, "s.whatsapp.net")
			services.SendReply(ctx, client, jid, msg)
		}
	})

	// 📩 Escuta eventos do WhatsApp
	events.Listen(ctx, client)

	// ⛔ Espera sinal de encerramento (Ctrl+C ou SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("📴 Encerrando conexão com o WhatsApp...")
	client.Disconnect()
}
