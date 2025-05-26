package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("🚀 Iniciando o bot WhatsApp...")

	ctx := context.Background()
	logger := waLog.Stdout("WhatsAppBot", "INFO", true)

	// Cria a conexão com o banco de dados SQLite
	db, err := sql.Open("sqlite", "file:session.db?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("❌ Erro ao abrir o banco de dados: %v", err)
	}

	// Configura os parâmetros de conexão para evitar SQLITE_BUSY
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Cria o container com a conexão configurada
	container := sqlstore.NewWithDB(db, "sqlite", logger)

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Println("⚠️ Nenhuma sessão ativa. Iniciando pareamento via QR Code...")
		deviceStore = container.NewDevice()
	}

	client := whatsmeow.NewClient(deviceStore, logger)

	if err := connectClient(ctx, client); err != nil {
		log.Fatalf("❌ Erro na conexão: %v", err)
	}

	log.Println("✅ Bot conectado ao WhatsApp!")

	// Escutando mensagens recebidas
	client.AddEventHandler(func(evt interface{}) {
		switch msg := evt.(type) {
		case *events.Message:
			if msg.Info.MessageSource.IsFromMe {
				return
			}

			text := msg.Message.GetConversation()
			if text == "" {
				return
			}

			sender := msg.Info.Sender.User
			log.Printf("📨 [%s]: %s", sender, text)

			handleCommand(ctx, client, msg.Info.Chat, text)
		}
	})

	// Aguarda sinal de encerramento (CTRL+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("⛔ Encerrando conexão...")
	client.Disconnect()
}

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
			return fmt.Errorf("❌ Erro ao escanear o QR")
		}
	}
	return nil
}

func handleCommand(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID, text string) {
	text = strings.ToLower(strings.TrimSpace(text))

	switch text {
	case "!ping":
		sendReply(ctx, client, chat, "pong ⚡")
	case "!help":
		sendReply(ctx, client, chat, "📖 Comandos disponíveis:\n!ping\n!help")
	default:
		if strings.Contains(text, "bom dia") {
			sendReply(ctx, client, chat, "🌞 Bom dia, guerreiro!")
		}
	}
}

func sendReply(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID, content string) {
	msg := &waProto.Message{
		Conversation: &content,
	}
	_, err := client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{})
	if err != nil {
		log.Printf("❌ Falha ao responder: %v", err)
	}
}
