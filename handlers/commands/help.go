package commands

import (
	"context"

	"github.com/faysk/whatsapp-bot/services"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

// Help mostra os comandos disponíveis
func Help(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID) {
	msg := `📖 *Comandos disponíveis*:
	
✅ !ping – Testa se o bot está online  
✅ !help – Mostra essa ajuda  
💬 Envie “bom dia” – Resposta especial`

	services.SendReply(ctx, client, chat, msg)
}
