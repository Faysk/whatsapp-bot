package services

import (
	"context"
	"log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

// SendReply envia uma mensagem de texto simples para um contato ou grupo
func SendReply(ctx context.Context, client *whatsmeow.Client, chat types.JID, content string) {
	if content == "" {
		log.Println("⚠️ Conteúdo vazio — mensagem não enviada.")
		return
	}

	msg := &proto.Message{
		Conversation: &content,
	}

	_, err := client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{})
	if err != nil {
		log.Printf("❌ Falha ao enviar mensagem para %s: %v", chat.String(), err)
	} else {
		log.Printf("📤 Mensagem enviada para %s", chat.String())
	}
}
