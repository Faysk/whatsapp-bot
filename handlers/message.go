package handlers

import (
	"context"
	"strings"

	"github.com/faysk/whatsapp-bot/handlers/commands"
	"github.com/faysk/whatsapp-bot/services"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

// HandleCommand interpreta o texto e direciona para o comando correspondente
func HandleCommand(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID, text string) {
	text = strings.ToLower(strings.TrimSpace(text))

	switch text {
	case "!ping":
		commands.Ping(ctx, client, chat)

	case "!help":
		commands.Help(ctx, client, chat)

	default:
		if strings.Contains(text, "bom dia") {
			commands.BomDia(ctx, client, chat)
		} else {
			// Pode responder mensagens genéricas aqui ou ignorar
			services.SendReply(ctx, client, chat, "🤖 Comando não reconhecido. Envie !help")
		}
	}
}
