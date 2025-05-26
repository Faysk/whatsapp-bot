package commands

import (
	"context"

	"github.com/faysk/whatsapp-bot/services"

	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

// BomDia responde a mensagens contendo "bom dia"
func BomDia(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID) {
	services.SendReply(ctx, client, chat, "🌞 Bom dia, guerreiro!")
}
