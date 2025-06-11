package events

import (
	"context"
	"log"

	"github.com/faysk/whatsapp-bot/handlers"
	"go.mau.fi/whatsmeow"
	waEvents "go.mau.fi/whatsmeow/types/events"
)

// Listen registra os listeners de eventos no cliente WhatsApp
func Listen(ctx context.Context, client *whatsmeow.Client) {
	client.AddEventHandler(func(evt interface{}) {
		switch msg := evt.(type) {

		case *waEvents.Message:
			// Ignora mensagens enviadas pelo próprio bot
			if msg.Info.MessageSource.IsFromMe {
				return
			}

			text := extractMessageText(msg)
			if text == "" {
				log.Printf("📭 Ignorando mensagem vazia ou não suportada de %s", msg.Info.Sender.User)
				return
			}

			log.Printf("📨 [%s] %s", msg.Info.Sender.User, text)
			handlers.HandleCommand(ctx, client, msg.Info.Chat, text, msg)

		// Futuro: adicionar suporte a eventos como presença, status, etc.
		default:
			// log.Printf("📡 Evento ignorado: %T", evt)
		}
	})
}

// extractMessageText extrai o conteúdo textual da mensagem recebida
func extractMessageText(msg *waEvents.Message) string {
	if msg.Message == nil {
		return ""
	}

	switch {
	case msg.Message.Conversation != nil:
		return *msg.Message.Conversation

	case msg.Message.ExtendedTextMessage != nil:
		return *msg.Message.ExtendedTextMessage.Text

	case msg.Message.ImageMessage != nil && msg.Message.ImageMessage.Caption != nil:
		return *msg.Message.ImageMessage.Caption

	case msg.Message.VideoMessage != nil && msg.Message.VideoMessage.Caption != nil:
		return *msg.Message.VideoMessage.Caption

	default:
		return ""
	}
}
