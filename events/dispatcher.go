package events

import (
	"context"
	"log"

	"github.com/faysk/whatsapp-bot/handlers"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// Listen registra os listeners de eventos no cliente WhatsApp
func Listen(ctx context.Context, client *whatsmeow.Client) {
	client.AddEventHandler(func(evt interface{}) {
		switch msg := evt.(type) {

		case *events.Message:
			if msg.Info.MessageSource.IsFromMe {
				return // ignora mensagens enviadas pelo próprio bot
			}

			text := extractMessageText(msg)
			if text == "" {
				log.Printf("📭 Ignorando mensagem vazia ou não suportada.")
				return
			}

			log.Printf("📨 [%s] %s", msg.Info.Sender.User, text)
			handlers.HandleCommand(ctx, client, msg.Info.Chat, text, msg)

		// Futuro: adicionar outros tipos de evento (ex: presence, status)
		default:
			// log.Printf("📡 Evento ignorado: %T", evt)
		}
	})
}

// extractMessageText extrai o conteúdo textual da mensagem recebida
func extractMessageText(msg *events.Message) string {
	if msg.Message == nil {
		return ""
	}

	// Tipos comuns de texto
	if msg.Message.Conversation != nil {
		return *msg.Message.Conversation
	}
	if msg.Message.ExtendedTextMessage != nil {
		return *msg.Message.ExtendedTextMessage.Text
	}

	// Outros tipos (opcional: áudio, imagem com legenda, etc.)
	return ""
}
