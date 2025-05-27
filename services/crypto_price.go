package services

import (
	"encoding/json"
	"fmt"
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
	for alias, id := range PredefinedAliases {
		cryptoAliases[alias] = id
	}

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

		if _, ok := cryptoAliases[id]; !ok {
			cryptoAliases[id] = id
		}
		if _, ok := cryptoAliases[symbol]; !ok {
			cryptoAliases[symbol] = id
		}
		if _, ok := cryptoAliases[name]; !ok {
			cryptoAliases[name] = id
		}
		noSpace := strings.ReplaceAll(name, " ", "")
		if _, ok := cryptoAliases[noSpace]; !ok {
			cryptoAliases[noSpace] = id
		}
	}
}

func GetCryptoPrice(input string) (string, error) {
	alias := strings.ToLower(strings.TrimSpace(input))
	cryptoID, ok := cryptoAliases[alias]
	if !ok {
		return "", fmt.Errorf("❌ Criptomoeda '%s' não reconhecida", input)
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false", cryptoID)

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("🌐 Erro ao acessar CoinGecko: %w", err)
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
		return "", fmt.Errorf("📦 Erro ao processar resposta: %w", err)
	}

	priceBRL := formatNumberBR(data.MarketData.CurrentPrice["brl"])
	priceUSD := formatNumberUS(data.MarketData.CurrentPrice["usd"])
	marketCap := formatNumberBR(data.MarketData.MarketCap["brl"])
	volume := formatNumberBR(data.MarketData.TotalVolume["brl"])
	rank := data.MarketData.MarketCapRank

	formatVar := func(val float64) string {
		switch {
		case val > 0:
			return fmt.Sprintf("🟢  %.2f%%", val)
		case val < 0:
			return fmt.Sprintf("🔴  %.2f%%", val)
		default:
			return fmt.Sprintf("⚪  %.2f%%", val)
		}
	}

	return fmt.Sprintf(
		"🪙 *%s (%s)*  |  🏅 Rank: #%d\n\n"+
			"💵 *Preço Atual*\n"+
			"🇧🇷 R$ %s\n"+
			"🇺🇸 $ %s\n\n"+
			"📊 *Variação*\n"+
			"1h:	%s\n"+
			"24h:	%s\n"+
			"7d:	%s\n"+
			"30d:	%s\n"+
			"1y:	%s\n\n"+
			"💰 *Market Cap:* R$ %s\n"+
			"📈 *Volume 24h:* R$ %s",
		data.Name,
		strings.ToUpper(data.Symbol),
		rank,
		priceBRL,
		priceUSD,
		formatVar(data.MarketData.PriceChangePercentage1h["brl"]),
		formatVar(data.MarketData.PriceChangePercentage24h),
		formatVar(data.MarketData.PriceChangePercentage7d),
		formatVar(data.MarketData.PriceChangePercentage30d),
		formatVar(data.MarketData.PriceChangePercentage1y),
		marketCap,
		volume,
	), nil
}
