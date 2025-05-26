package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type coinInfo struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type CoinData struct {
	ID     string
	Name   string
	Symbol string
}

var (
	cryptoAliases = make(map[string]string)   // alias → id
	cryptoInfoMap = make(map[string]CoinData) // id → CoinData
	client        = &http.Client{Timeout: 10 * time.Second}
)

// init carrega as moedas e gera os mapas de aliases automaticamente
func init() {
	resp, err := client.Get("https://api.coingecko.com/api/v3/coins/list")
	if err != nil {
		fmt.Println("⚠️ Erro ao consultar CoinGecko:", err)
		return
	}
	defer resp.Body.Close()

	var coins []coinInfo
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		fmt.Println("⚠️ Erro ao decodificar moedas:", err)
		return
	}

	for _, coin := range coins {
		id := coin.ID
		name := strings.ToLower(coin.Name)
		symbol := strings.ToLower(coin.Symbol)

		cryptoInfoMap[id] = CoinData{
			ID:     id,
			Name:   coin.Name,
			Symbol: strings.ToUpper(coin.Symbol),
		}

		// Vários aliases para facilitar comandos
		cryptoAliases[id] = id
		cryptoAliases[symbol] = id
		cryptoAliases[name] = id
		cryptoAliases[strings.ReplaceAll(name, " ", "")] = id
	}
}

// GetCryptoPrice retorna o preço formatado de qualquer cripto em BRL e USD
func GetCryptoPrice(input string) (string, error) {
	alias := strings.ToLower(strings.TrimSpace(input))
	cryptoID, ok := cryptoAliases[alias]
	if !ok {
		return "", fmt.Errorf("❌ Criptomoeda '%s' não reconhecida", input)
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=brl,usd", cryptoID)
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("🌐 Erro ao acessar CoinGecko: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("❌ CoinGecko retornou status %d", resp.StatusCode)
	}

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("📦 Erro ao processar dados: %w", err)
	}

	prices, ok := data[cryptoID]
	if !ok {
		return "", fmt.Errorf("⚠️ Criptomoeda '%s' não encontrada", cryptoID)
	}

	coin := cryptoInfoMap[cryptoID]

	return fmt.Sprintf(
		"💰 *%s (%s)*\n🇧🇷 R$ %s\n🇺🇸 $ %s",
		coin.Name,
		coin.Symbol,
		formatFloat(prices["brl"]),
		formatFloat(prices["usd"]),
	), nil
}

// formatFloat retorna valor com separador decimal padronizado
func formatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 2, 64)
}
