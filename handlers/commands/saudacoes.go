package commands

import (
	"context"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/faysk/whatsapp-bot/services"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

var saudacoes = map[string][]string{
	"bom dia": {
		"🌞 Bom dia, guerreiro!",
		"☕ Bom dia! Que hoje seja incrível!",
		"👊 Bom dia, bora vencer mais um dia!",
		"✨ Bom dia! Nova chance de vencer.",
	},
	"boa tarde": {
		"🌤️ Boa tarde! Força total no meio do dia!",
		"📈 Boa tarde! Produtividade nas alturas!",
		"⚡ Boa tarde! Vamos com tudo!",
	},
	"boa noite": {
		"🌙 Boa noite! Hora de recarregar.",
		"🛌 Boa noite! Durma com os anjos.",
		"😴 Boa noite, guerreiro. Amanhã tem mais luta!",
	},
	"oi": {
		"👋 Oi! Tudo certo por aí?",
		"E aí! Como vai você?",
		"Salve, salve!",
		"Opa! Cheguei na hora certa?",
	},
	"olá": {
		"Olá! Seja bem-vindo!",
		"Oiê! Chegou bem na hora.",
		"👋 Olá! Tudo tranquilo?",
	},
	"salve": {
		"👊 Salve, parceiro!",
		"Salve! Tamo junto!",
		"Salve, salve! Que a força esteja com você.",
	},
	"opa": {
		"Opa! Tudo beleza?",
		"E aí, opa!",
		"Fala aí! 👋",
	},
}

var aliases = map[string][]string{
	"bom dia":    {"bom dia", "b dia", "bdia"},
	"boa tarde":  {"boa tarde", "boa tardi", "boa trde"},
	"boa noite":  {"boa noite", "boanoite", "boa noti"},
	"oi":         {"oi", "e aí", "eai", "iae", "oii"},
	"olá":        {"olá", "ola", "olaaa"},
	"salve":      {"salve", "salvee", "salvee!"},
	"opa":        {"opa", "opaaa", "oopaa"},
}

// DetectSaudacao tenta identificar e responder a uma saudação
func DetectSaudacao(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID, text string) bool {
	text = normalize(text)
	words := tokenize(text)
	rand.Seed(time.Now().UnixNano())

	for categoria, variações := range aliases {
		for _, termo := range variações {
			for _, palavra := range words {
				if palavra == termo {
					respostas := saudacoes[categoria]
					if len(respostas) > 0 {
						resposta := respostas[rand.Intn(len(respostas))]
						services.SendReply(ctx, client, chat, resposta)
						return true
					}
				}
			}
		}
	}

	r
