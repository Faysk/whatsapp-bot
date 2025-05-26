package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetBitcoinPrice consulta a cotação atual do Bitcoin em BRL via CoinGecko
func GetBitcoinPrice() (string, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=brl"

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("erro ao acessar CoinGecko: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("falha na API CoinGecko: status %d", resp.StatusCode)
	}

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta da API: %w", err)
	}

	price, ok := data["bitcoin"]["brl"]
	if !ok {
		return "", fmt.Errorf("não foi possível encontrar o preço do Bitcoin")
	}

	return fmt.Sprintf("💰 O preço atual do *Bitcoin* é *R$ %.2f*", price), nil
}
