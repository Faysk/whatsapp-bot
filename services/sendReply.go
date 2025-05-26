package services

import (
	"context"
	"log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	p "google.golang.org/protobuf/proto"
)

// SendReply envia uma mensagem para um JID (grupo ou contato)
func SendReply(ctx context.Context, client *whatsmeow.Client, chat types.JID, content string) {
	if content == "" {
		log.Println("⚠️ Conteúdo vazio — mensagem não enviada.")
		return
	}

	msg := &proto.Message{
		Conversation: p.String(content),
	}

	_, err := client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{})
	if err != nil {
		log.Printf("❌ Falha ao enviar mensagem para %s: %v", chat.String(), err)
	} else {
		log.Printf("📤 Mensagem enviada para %s", chat.String())
	}
}

// SendToNumber envia mensagem diretamente para um número com formato 5511...
func SendToNumber(ctx context.Context, client *whatsmeow.Client, phone string, content string) {
	if phone == "" || content == "" {
		log.Println("⚠️ Número ou conteúdo vazio — mensagem não enviada.")
		return
	}

	jid := types.NewJID(phone, "s.whatsapp.net")
	msg := &proto.Message{
		Conversation: p.String(content),
	}

	_, err := client.SendMessage(ctx, jid, msg, whatsmeow.SendRequestExtra{})
	if err != nil {
		log.Printf("❌ Erro ao enviar mensagem para %s: %v", phone, err)
	} else {
		log.Printf("📬 Mensagem enviada para %s", phone)
	}
}
