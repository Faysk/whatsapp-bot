package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/faysk/whatsapp-bot/config"
	"github.com/faysk/whatsapp-bot/events"
	"github.com/faysk/whatsapp-bot/store"
	"github.com/faysk/whatsapp-bot/utils"
)

func main() {
	log.Println("🚀 Iniciando o bot WhatsApp...")

	// Inicializa configurações (ex: leitura de .env futuramente)
	config.Load()

	// Inicializa logger com cor e timestamps, se necessário
	utils.SetupLogger()

	// Cria contexto de execução principal
	ctx := context.Background()

	// Inicializa o cliente WhatsApp com persistência
	client, err := store.InitClient(ctx)
	if err != nil {
		log.Fatalf("❌ Erro ao iniciar o cliente WhatsApp: %v", err)
	}

	// Escuta e despacha os eventos recebidos
	events.Listen(ctx, client)

	// Aguarda sinal de encerramento (CTRL+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("⛔ Encerrando conexão...")
	client.Disconnect()
}
