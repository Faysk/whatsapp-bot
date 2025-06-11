package commands

import (
	"context"
	"github.com/faysk/whatsapp-bot/services"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

// Help mostra os comandos e interações disponíveis com o bot
func Help(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID) {
	msg := `📖 *Comandos e interações disponíveis*:

🧪 *Comandos tradicionais*:
- !ping → Testa se o bot está online
- !help → Exibe esta mensagem de ajuda

🤖 *Interações naturais com o bot*:
- Diga: "ping", "teste", "tá aí", "responde", etc.
- O bot vai responder com frases aleatórias

🌞 *Saudações automáticas*:
- "bom dia", "boa tarde", "boa noite"
- "oi", "olá", "salve", "opa"

ℹ️ *Mais funções em breve...*

💡 Dica: use linguagem natural! O bot entende mais do que apenas comandos. 😉`

	services.SendReply(ctx, client, chat, msg)
}
