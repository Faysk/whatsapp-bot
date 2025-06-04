package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type CryptoRecord struct {
	AllTimeHigh float64   `json:"ath"`
	Timestamp   time.Time `json:"timestamp"`
}

var monitoredCoins = []string{"BTC", "ETH", "USDT", "XRP", "SOL"}

const (
	priceAPIURL    = "https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd"
	recordFilePath = "crypto_records.json"
	checkInterval  = 5 * time.Minute
)

func MonitorCryptos(sendAlert func(string)) {
	ticker := time.NewTicker(checkInterval)

	go func() {
		for range ticker.C {
			log.Println("🔍 Verificando máximas históricas (ATH oficiais)...")

			for _, symbol := range monitoredCoins {
				id := getCoinGeckoID(symbol)
				url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false", id)

				resp, err := http.Get(url)
				if err != nil {
					log.Printf("❌ %s: erro HTTP ao acessar CoinGecko: %v", symbol, err)
					continue
				}

				if resp.StatusCode != http.StatusOK {
					log.Printf("❌ %s: CoinGecko retornou status %d", symbol, resp.StatusCode)
					resp.Body.Close()
					continue
				}

				var data struct {
					MarketData struct {
						CurrentPrice map[string]float64 `json:"current_price"`
						ATH          map[string]float64 `json:"ath"`
					} `json:"market_data"`
				}

				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					log.Printf("❌ %s: erro ao decodificar resposta: %v", symbol, err)
					resp.Body.Close()
					continue
				}
				resp.Body.Close()

				current := data.MarketData.CurrentPrice["usd"]
				ath := data.MarketData.ATH["usd"]

				if current > ath {
					log.Printf("🚀 %s ultrapassou ATH oficial: U$ %.2f > ATH U$ %.2f", symbol, current, ath)

					detail, err := GetCryptoPriceWithOverride(symbol, current)
					if err != nil {
						log.Printf("❌ %s: erro ao montar mensagem: %v", symbol, err)
						continue
					}

					alert := fmt.Sprintf("🚨 *NOVO RECORD HISTÓRICO (ATH)*\n\n%s\n\n🕒 ATH oficial superado em %s", detail, time.Now().Format("02/01/2006 15:04"))
					sendAlert(alert)
				} else {
					log.Printf("ℹ️ %s: atual U$ %.2f — ainda abaixo do ATH oficial U$ %.2f", symbol, current, ath)
				}
			}
		}
	}()
}

func getCryptoPrice(coinID string) (float64, error) {
	url := fmt.Sprintf(priceAPIURL, coinID)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("erro HTTP: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("erro JSON: %w", err)
	}

	price, ok := result[coinID]["usd"]
	if !ok {
		return 0, fmt.Errorf("preço ausente para %s", coinID)
	}
	return price, nil
}

func loadRecords() map[string]CryptoRecord {
	data, err := os.ReadFile(recordFilePath)
	if err != nil {
		log.Printf("📂 Nenhum recorde salvo, iniciando vazio...")
		return make(map[string]CryptoRecord)
	}
	var records map[string]CryptoRecord
	if err := json.Unmarshal(data, &records); err != nil {
		log.Printf("⚠️ Erro ao ler JSON: %v", err)
		return make(map[string]CryptoRecord)
	}
	return records
}

func saveRecords(records map[string]CryptoRecord) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Printf("❌ Erro ao codificar JSON: %v", err)
		return
	}
	if err := os.WriteFile(recordFilePath, data, 0644); err != nil {
		log.Printf("❌ Erro ao salvar JSON: %v", err)
	}
}

func getCoinGeckoID(symbol string) string {
	switch strings.ToUpper(symbol) {
	case "BTC":
		return "bitcoin"
	case "ETH":
		return "ethereum"
	case "USDT":
		return "tether"
	case "XRP":
		return "ripple"
	case "SOL":
		return "solana"
	default:
		return symbol
	}
}

// 🔽 NOVA FUNÇÃO - sem interferir na original usada pelo !btc, !eth etc.
func GetCryptoPriceWithOverride(input string, usdOverride float64) (string, error) {
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
	priceUSD := formatNumberUS(usdOverride)
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
