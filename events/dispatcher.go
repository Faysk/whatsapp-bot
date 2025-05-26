package events

import (
	"context"
	"log"

	"github.com/faysk/whatsapp-bot/handlers"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// Listen registra os listeners de eventos no client WhatsApp
func Listen(ctx context.Context, client *whatsmeow.Client) {
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

			log.Printf("📨 [%s]: %s", msg.Info.Sender.User, text)
			handlers.HandleCommand(ctx, client, msg.Info.Chat, text)

		default:
			// Aqui você pode tratar outros tipos de evento no futuro
		}
	})
}
