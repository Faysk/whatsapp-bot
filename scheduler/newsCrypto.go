package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/faysk/whatsapp-bot/services"
	"github.com/go-co-op/gocron"
	"go.mau.fi/whatsmeow"
)

// StartDailyNews agenda o envio diário de notícias de criptomoedas às 10h (horário local)
func StartDailyNews(ctx context.Context, client *whatsmeow.Client, numbers []string) {
	if len(numbers) == 0 {
		log.Println("⚠️ Nenhum número autorizado para envio de notícias.")
		return
	}

	s := gocron.NewScheduler(time.Local)

	_, err := s.Every(1).Day().At("10:00").Tag("daily-crypto-news").Do(func() {
		sendCryptoNews(ctx, client, numbers)
	})

	if err != nil {
		log.Printf("❌ Erro ao agendar job de notícias cripto: %v", err)
		return
	}

	log.Println("🗞️ Agendador de notícias cripto ativado (todos os dias às 10h)")
	s.StartAsync()
}

// sendCryptoNews busca e envia as últimas atualizações de criptomoedas em dois blocos (Trending + News)
func sendCryptoNews(ctx context.Context, client *whatsmeow.Client, numbers []string) {
	log.Println("📡 Iniciando coleta de notícias do CryptoPanic...")

	trendingMsg, newsMsg, err := services.GetCryptoNews()
	if err != nil {
		log.Printf("❌ Erro ao obter notícias: %v", err)
		return
	}

	if trendingMsg == "" && newsMsg == "" {
		log.Println("⚠️ Nenhuma notícia relevante disponível no momento.")
		return
	}

	log.Printf("📦 Notícias prontas para envio: Trending (%d caracteres), News (%d caracteres)",
		len(trendingMsg), len(newsMsg),
	)

	for _, number := range numbers {
		if number == "" {
			continue
		}

		if trendingMsg != "" {
			log.Printf("📤 Enviando 🔥 *Tópicos em Alta* para %s", number)
			services.SendToNumber(ctx, client, number, trendingMsg)
			time.Sleep(2 * time.Second) // Delay para evitar truncamento
		}

		if newsMsg != "" {
			log.Printf("📤 Enviando 🗞️ *Últimas Notícias* para %s", number)
			services.SendToNumber(ctx, client, number, newsMsg)
			time.Sleep(2 * time.Second)
		}
	}

	log.Println("✅ Todas as notícias cripto foram enviadas com sucesso.")
}
