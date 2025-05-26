package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/faysk/whatsapp-bot/config"
	"github.com/faysk/whatsapp-bot/events"
	"github.com/faysk/whatsapp-bot/services"
	"github.com/faysk/whatsapp-bot/store"
	"github.com/faysk/whatsapp-bot/utils"
)

func main() {
	log.Println("🚀 Iniciando o bot WhatsApp...")

	// 🔧 Carrega configurações do .env e variáveis globais
	config.Load()

	// 🔐 Carrega números autorizados dinâmicos do arquivo JSON
	dynamic := store.LoadAuthorizedNumbers()
	config.AddDynamicAuthorizedNumbers(dynamic)

	// 📝 Inicializa logger (timestamp + cores)
	utils.SetupLogger()

	// 🌐 Cria contexto de execução principal
	ctx := context.Background()

	// 📲 Inicializa cliente WhatsApp com persistência
	client, err := services.InitWhatsAppClient(ctx)
	if err != nil {
		log.Fatalf("❌ Falha ao inicializar o cliente WhatsApp: %v", err)
	}

	// 🔗 Conecta automaticamente com sessão ou mostra QR Code
	if err := services.ConnectWithQR(ctx, client); err != nil {
		log.Fatalf("❌ Erro na conexão com WhatsApp: %v", err)
	}

	log.Println("✅ Bot conectado com sucesso. Aguardando mensagens...")

	// 📩 Inicia o listener de eventos (mensagens, grupos, etc)
	events.Listen(ctx, client)

	// 🛑 Espera até que o processo receba uma interrupção (Ctrl+C ou SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("⛔ Encerrando conexão com o WhatsApp...")
	client.Disconnect()
}
