package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/faysk/whatsapp-bot/openai" // necessário se quiser traduzir
)

const cryptoPanicToken = "444829389373d26117c89452f7bb83efbbc4524d"

type panicPost struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type panicResponse struct {
	Results []panicPost `json:"results"`
}

// GetCryptoNews retorna até 10 manchetes cripto mais recentes, traduzidas para PT-BR
func GetCryptoNews() (string, error) {
	url := fmt.Sprintf("https://cryptopanic.com/api/v1/posts/?auth_token=%s&kind=news&public=true", cryptoPanicToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("erro ao acessar CryptoPanic: %w", err)
	}
	defer resp.Body.Close()

	var data panicResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	if len(data.Results) == 0 {
		return "⚠️ Nenhuma notícia cripto encontrada hoje.", nil
	}

	builder := strings.Builder{}
	builder.WriteString("📰 *Notícias Cripto de Hoje:*\n\n")

	for i, post := range data.Results {
		if i >= 10 {
			break
		}
		builder.WriteString(fmt.Sprintf("• %s\n🔗 https://cryptopanic.com/news/%s\n\n", post.Title, post.Slug))
	}

	// Tradução automática via ChatGPT
	traduzido, err := openai.AskChatGPT("Traduza esse conteúdo para português mantendo o formato:\n\n" + builder.String())
	if err != nil {
		// fallback para inglês caso tradução falhe
		return builder.String() + "\n⚠️ Falha ao traduzir: " + err.Error(), nil
	}

	return traduzido, nil
}
