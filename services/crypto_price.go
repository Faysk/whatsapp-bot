// File: services/crypto_price.go
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type CoinData struct {
	ID     string
	Name   string
	Symbol string
}

type coinInfo struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

var (
	cryptoAliases = make(map[string]string)   // alias → id
	cryptoInfoMap = make(map[string]CoinData) // id → CoinData
	client        = &http.Client{Timeout: 10 * time.Second}
)

func init() {
	log.Println("🔄 Carregando aliases de criptomoedas...")

	// Adiciona aliases fixos
	for alias, id := range PredefinedAliases {
		cryptoAliases[strings.ToLower(alias)] = id
	}

	// Consulta CoinGecko
	resp, err := client.Get("https://api.coingecko.com/api/v3/coins/list")
	if err != nil {
		log.Printf("⚠️ Erro ao acessar CoinGecko: %v", err)
		return
	}
	defer resp.Body.Close()

	var coins []coinInfo
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		log.Printf("⚠️ Erro ao decodificar resposta do CoinGecko: %v", err)
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

		// Evita sobrescrever aliases manuais
		addAliasIfMissing(id, id)
		addAliasIfMissing(symbol, id)
		addAliasIfMissing(name, id)
		addAliasIfMissing(strings.ReplaceAll(name, " ", ""), id)
	}

	log.Printf("✅ %d criptomoedas carregadas de CoinGecko.", len(cryptoInfoMap))
}

func addAliasIfMissing(alias, id string) {
	if _, exists := cryptoAliases[alias]; !exists {
		cryptoAliases[alias] = id
	}
}

// GetCryptoPrice retorna a cotação formatada de uma moeda
func GetCryptoPrice(input string) (string, error) {
	alias := strings.ToLower(strings.TrimSpace(input))
	cryptoID, ok := cryptoAliases[alias]
	if !ok {
		return "", fmt.Errorf("❌ Criptomoeda '%s' não reconhecida", input)
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false", cryptoID)

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("🌐 Erro HTTP ao acessar CoinGecko: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("❌ CoinGecko retornou status %d", resp.StatusCode)
	}

	var data struct {
		Name       string `json:"name"`
		Symbol     string `json:"symbol"`
		MarketData struct {
			CurrentPrice             map[string]float64 `json:"current_price"`
			MarketCap                map[string]float64 `json:"market_cap"`
			TotalVolume              map[string]float64 `json:"total_volume"`
			MarketCapRank            int                `json:"market_cap_rank"`
			PriceChangePercentage1h  map[string]float64 `json:"price_change_percentage_1h_in_currency"`
			PriceChangePercentage24h float64            `json:"price_change_percentage_24h"`
			PriceChangePercentage7d  float64            `json:"price_change_percentage_7d"`
			PriceChangePercentage30d float64            `json:"price_change_percentage_30d"`
			PriceChangePercentage1y  float64            `json:"price_change_percentage_1y"`
		} `json:"market_data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("📦 Erro ao decodificar resposta: %w", err)
	}

	formatVar := func(val float64) string {
		switch {
		case val > 0:
			return fmt.Sprintf("🟢 %.2f%%", val)
		case val < 0:
			return fmt.Sprintf("🔴 %.2f%%", val)
		default:
			return fmt.Sprintf("⚪ %.2f%%", val)
		}
	}

	return fmt.Sprintf(
		"🪙 *%s (%s)*  |  🏅 Rank: #%d\n\n"+
			"💵 *Preço Atual*\n🇧🇷 R$ %s\n🇺🇸 $ %s\n\n"+
			"📊 *Variação*\n1h: %s\n24h: %s\n7d: %s\n30d: %s\n1y: %s\n\n"+
			"💰 *Market Cap:* R$ %s\n📈 *Volume 24h:* R$ %s",
		data.Name,
		strings.ToUpper(data.Symbol),
		data.MarketData.MarketCapRank,
		formatNumberBR(data.MarketData.CurrentPrice["brl"]),
		formatNumberUS(data.MarketData.CurrentPrice["usd"]),
		formatVar(data.MarketData.PriceChangePercentage1h["brl"]),
		formatVar(data.MarketData.PriceChangePercentage24h),
		formatVar(data.MarketData.PriceChangePercentage7d),
		formatVar(data.MarketData.PriceChangePercentage30d),
		formatVar(data.MarketData.PriceChangePercentage1y),
		formatNumberBR(data.MarketData.MarketCap["brl"]),
		formatNumberBR(data.MarketData.TotalVolume["brl"]),
	), nil
}
