package commands

import (
	"context"

	"github.com/faysk/whatsapp-bot/services"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

// Ping responde com "Pong!"
func Ping(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID) {
	services.SendReply(ctx, client, chat, "🏓 Pong!")
}
