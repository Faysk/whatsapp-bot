package commands

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/faysk/whatsapp-bot/services"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

var interacoes = map[string][]string{
	"ping": {
		"🏓 Pong!",
		"✅ Estou online!",
		"💡 Funcionando perfeitamente.",
		"📶 Sinal forte por aqui.",
	},
	"teste": {
		"🔍 Teste recebido com sucesso!",
		"✅ Testado e aprovado.",
		"🎯 Está tudo funcionando!",
	},
	"tá aí": {
		"🙋‍♂️ Estou aqui!",
		"👀 Sempre observando...",
		"🤖 Operacional e aguardando comandos.",
	},
	"bot": {
		"Sim, senhor! 🤖",
		"Chamou o bot? Cheguei!",
	},
}

// DetectInteracao verifica se a mensagem é uma interação comum com o bot
func DetectInteracao(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID, text string) bool {
	text = normalize(text)
	rand.Seed(time.Now().UnixNano())

	for chave, respostas := range interacoes {
		if strings.Contains(text, chave) {
			resposta := respostas[rand.Intn(len(respostas))]
			services.SendReply(ctx, client, chat, resposta)
			return true
		}
	}
	return false
}

// Funções auxiliares

// normalize simplifica acentuação e minúsculas
func normalize(input string) string {
	input = strings.ToLower(input)
	var result strings.Builder
	for _, r := range input {
		switch r {
		case 'á', 'à', 'ã', 'â':
			result.WriteRune('a')
		case 'é', 'ê':
			result.WriteRune('e')
		case 'í':
			result.WriteRune('i')
		case 'ó', 'õ', 'ô':
			result.WriteRune('o')
		case 'ú':
			result.WriteRune('u')
		case 'ç':
			result.WriteRune('c')
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}
