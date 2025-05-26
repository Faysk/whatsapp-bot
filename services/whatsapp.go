package services

import (
	"context"
	"log"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
)

// SendReply envia uma mensagem de texto simples para um contato/grupo
func SendReply(ctx context.Context, client *whatsmeow.Client, chat waTypes.JID, content string) {
	msg := &waProto.Message{
		Conversation: &content,
	}

	_, err := client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{})
	if err != nil {
		log.Printf("❌ Falha ao enviar mensagem: %v", err)
	}
}
