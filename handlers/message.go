package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/faysk/whatsapp-bot/config"
	"github.com/faysk/whatsapp-bot/handlers/commands"
	"github.com/faysk/whatsapp-bot/openai"
	"github.com/faysk/whatsapp-bot/services"
	"github.com/faysk/whatsapp-bot/store"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func HandleCommand(ctx context.Context, client *whatsmeow.Client, _ waTypes.JID, text string, msg *events.Message) {
	text = strings.TrimSpace(text)
	lower := strings.ToLower(text)

	logPrefix := fmt.Sprintf("[%s]", config.AppConfig.BotName)
	sender := msg.Info.Sender.User
	isGroup := msg.Info.IsGroup
	chat := msg.Info.Chat

	// 🔒 Restrições
	if config.AppConfig.RestrictToGroup && !isGroup {
		log.Printf("%s 🚫 Ignorando mensagem privada (RESTRICT_TO_GROUP=true)", logPrefix)
		return
	}
	if !isAuthorized(sender) {
		log.Printf("%s 🚫 Número não autorizado: %s", logPrefix, sender)
		return
	}

	// ✅ Comandos diretos
	switch lower {
	case "!ping":
		log.Printf("%s 🟢 Comando !ping de %s", logPrefix, sender)
		commands.Ping(ctx, client, chat)
		return

	case "!help":
		log.Printf("%s 📘 Comando !help de %s", logPrefix, sender)
		commands.Help(ctx, client, chat)
		return

	case "!cryptonews":
		log.Printf("%s 📰 Comando !cryptonews de %s", logPrefix, sender)
		news, err := services.GetCryptoNews()
		if err != nil || news == "" {
			msg := "⚠️ Não foi possível obter as notícias de criptomoedas no momento."
			if err != nil {
				msg += "\nDetalhes: " + err.Error()
			}
			services.SendReply(ctx, client, chat, msg)
			return
		}
		services.SendReply(ctx, client, chat, news)
		return
	}

	// 💱 Comando de Cripto (qualquer !moeda, ex: !btc, !sol)
	if strings.HasPrefix(lower, "!") {
		moeda := strings.TrimPrefix(lower, "!")
		log.Printf("%s 💰 Comando de cripto '%s' de %s", logPrefix, moeda, sender)

		price, err := services.GetCryptoPrice(moeda)
		if err != nil {
			price = "❌ Erro ao consultar moeda: " + err.Error()
		}
		services.SendReply(ctx, client, chat, price)
		return
	}

	// 🌞 Saudações simples
	if strings.Contains(lower, "bom dia") {
		log.Printf("%s ☀️ Saudação detectada de %s", logPrefix, sender)
		commands.BomDia(ctx, client, chat)
		return
	}

	// 🤖 IA: qualquer frase com "renan"
	if config.AppConfig.EnableChatGPT && config.AppConfig.OpenAIKey != "" && strings.Contains(lower, "renan") {

		if strings.Contains(lower, "adicione o número") {
			num := extractPhoneNumber(lower)
			if num == "" {
				services.SendReply(ctx, client, chat, "⚠️ Nenhum número válido encontrado.")
				return
			}
			if err := store.AddAuthorized(num); err != nil {
				services.SendReply(ctx, client, chat, "⚠️ Não foi possível adicionar o número.")
				return
			}
			config.AddDynamicAuthorizedNumbers([]string{num})
			log.Printf("%s ➕ Número %s adicionado por %s", logPrefix, num, sender)
			services.SendReply(ctx, client, chat, fmt.Sprintf("✅ Número %s adicionado à lista de autorizados.", num))
			return
		}

		if strings.Contains(lower, "remova o número") {
			num := extractPhoneNumber(lower)
			if num == "" {
				services.SendReply(ctx, client, chat, "⚠️ Nenhum número válido encontrado.")
				return
			}
			if err := store.RemoveAuthorized(sender, num); err != nil {
				services.SendReply(ctx, client, chat, fmt.Sprintf("⚠️ %v", err))
				return
			}
			config.AppConfig.AuthorizedNumbers = store.LoadAuthorizedNumbers()
			log.Printf("%s ➖ Número %s removido por %s", logPrefix, num, sender)
			services.SendReply(ctx, client, chat, fmt.Sprintf("🗑️ Número %s removido da lista de autorizados.", num))
			return
		}

		log.Printf("%s 🤖 Enviando mensagem para IA: \"%s\" de %s", logPrefix, text, sender)
		reply, err := openai.AskChatGPT(text)
		if err != nil {
			log.Printf("%s ⚠️ Erro na IA: %v", logPrefix, err)
			reply = "❌ Erro ao consultar a IA: " + err.Error()
		}
		services.SendReply(ctx, client, chat, reply)
		return
	}

	// 🔇 Ignorado
	log.Printf("%s ❌ Ignorado: \"%s\" de %s (sem comando nem palavra-chave)", logPrefix, text, sender)
}

// isAuthorized verifica se o número está autorizado
func isAuthorized(sender string) bool {
	for _, num := range config.AppConfig.AuthorizedNumbers {
		if sender == num {
			return true
		}
	}
	return false
}

// extractPhoneNumber tenta extrair o primeiro número válido
func extractPhoneNumber(text string) string {
	words := strings.Fields(text)
	for _, word := range words {
		clean := strings.Trim(word, ",. ")
		if strings.HasPrefix(clean, "55") && len(clean) >= 11 {
			return clean
		}
	}
	return ""
}
