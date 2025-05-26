package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/faysk/whatsapp-bot/services"
	"github.com/go-co-op/gocron"
	"go.mau.fi/whatsmeow"
)

// StartDailyNews agenda o envio automático de notícias de criptomoedas às 10h da manhã (horário local)
func StartDailyNews(ctx context.Context, client *whatsmeow.Client, numbers []string) {
	s := gocron.NewScheduler(time.Local)

	s.Every(1).Day().At("10:00").Do(func() {
		sendCryptoNews(ctx, client, numbers)
	})

	log.Println("🗞️ Agendador de notícias cripto ativado (todos os dias às 10h)")
	s.StartAsync()
}

// sendCryptoNews busca e envia as notícias mais recentes para os números autorizados
func sendCryptoNews(ctx context.Context, client *whatsmeow.Client, numbers []string) {
	log.Println("📡 Buscando notícias de cripto para envio...")

	news, err := services.GetCryptoNews()
	if err != nil {
		log.Printf("❌ Erro ao buscar notícias: %v", err)
		return
	}

	if news == "" {
		log.Println("⚠️ Nenhuma notícia relevante disponível no momento.")
		return
	}

	for _, number := range numbers {
		if number == "" {
			continue
		}
		services.SendToNumber(ctx, client, number, news)
		log.Printf("📬 Notícia enviada para %s", number)
	}
}
