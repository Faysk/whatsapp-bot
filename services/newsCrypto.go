// File: services/newsCrypto.go
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/faysk/whatsapp-bot/openai"
)

const (
	cryptoPanicToken   = "444829389373d26117c89452f7bb83efbbc4524d"
	cryptoPanicBaseURL = "https://cryptopanic.com/news/"
	maxItemsPerSection = 10
	newsRequestTimeout = 10 * time.Second
)

var newsClient = &http.Client{Timeout: newsRequestTimeout}

type panicPost struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type panicResponse struct {
	Results []panicPost `json:"results"`
}

func fetchCryptoPanic(endpoint string) ([]panicPost, error) {
	resp, err := newsClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar CryptoPanic: %w", err)
	}
	defer resp.Body.Close()

	var data panicResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return data.Results, nil
}

func removeDuplicates(trending, news []panicPost) []panicPost {
	seen := make(map[string]struct{})
	for _, post := range trending {
		seen[post.Slug] = struct{}{}
	}

	filtered := make([]panicPost, 0, len(news))
	for _, post := range news {
		if _, exists := seen[post.Slug]; !exists {
			filtered = append(filtered, post)
		}
	}
	return filtered
}

func GetCryptoNews() (string, string, error) {
	hotURL := fmt.Sprintf("https://cryptopanic.com/api/v1/posts/?auth_token=%s&filter=hot&public=true", cryptoPanicToken)
	allURL := fmt.Sprintf("https://cryptopanic.com/api/v1/posts/?auth_token=%s&kind=news&public=true", cryptoPanicToken)

	hotPosts, err := fetchCryptoPanic(hotURL)
	if err != nil {
		return "", "", err
	}

	newsPosts, err := fetchCryptoPanic(allURL)
	if err != nil {
		return "", "", err
	}

	newsPosts = removeDuplicates(hotPosts, newsPosts)

	var hotBuilder strings.Builder
	hotBuilder.WriteString("📰 *Resumo Cripto do Dia*\n")
	hotBuilder.WriteString("──────────────────────\n\n🔥 *Mais Quentes do Momento*\n\n")
	if len(hotPosts) == 0 {
		hotBuilder.WriteString("⚠️ Nenhuma notícia quente no momento.\n")
	} else {
		for i, post := range hotPosts {
			if i >= maxItemsPerSection {
				break
			}
			hotBuilder.WriteString(fmt.Sprintf("*%d.* %s\n\n", i+1, strings.TrimSpace(post.Title)))
		}
		hotBuilder.WriteString(fmt.Sprintf("🔗 (Fonte: %s)\n", cryptoPanicBaseURL))
	}

	var newsBuilder strings.Builder
	newsBuilder.WriteString("🗞️ *Últimas Notícias*\n")
	newsBuilder.WriteString("──────────────────────\n\n")
	if len(newsPosts) == 0 {
		newsBuilder.WriteString("⚠️ Nenhuma notícia recente disponível.\n")
	} else {
		for i, post := range newsPosts {
			if i >= maxItemsPerSection {
				break
			}
			newsBuilder.WriteString(fmt.Sprintf("*%d.* %s\n\n", i+1, strings.TrimSpace(post.Title)))
		}
		newsBuilder.WriteString(fmt.Sprintf("🔗 (Fonte: %s)\n", cryptoPanicBaseURL))
	}

	translate := func(txt string) string {
		result, err := openai.AskChatGPT(
			"Traduza para português mantendo estilo limpo, direto e compatível com WhatsApp. Mantenha a estrutura e link no final:\n\n" + txt,
		)
		if err != nil || strings.TrimSpace(result) == "" {
			log.Println("⚠️ Falha na tradução:", err)
			return txt
		}
		return result
	}

	translatedHot := translate(hotBuilder.String())
	translatedNews := translate(newsBuilder.String())

	log.Printf("🔍 Total hot: %d | Total news (filtradas): %d", len(hotPosts), len(newsPosts))

	return strings.TrimSpace(translatedHot), strings.TrimSpace(translatedNews), nil
}
