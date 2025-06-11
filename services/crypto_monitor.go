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
	coinDetailAPI  = "https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false"
	recordFilePath = "crypto_records.json"
	checkInterval  = 5 * time.Minute
	requestTimeout = 10 * time.Second
)

func MonitorCryptos(sendAlert func(string)) {
	ticker := time.NewTicker(checkInterval)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("🔥 Panic recuperado no monitor de criptos: %v", r)
			}
		}()

		records := loadRecords()
		client := &http.Client{Timeout: requestTimeout}

		for range ticker.C {
			log.Println("🔍 Verificando máximas históricas (ATH oficiais)...")

			for _, symbol := range monitoredCoins {
				id := ResolveAlias(symbol)
				url := fmt.Sprintf(coinDetailAPI, id)

				resp, err := client.Get(url)
				if err != nil {
					log.Printf("❌ [%s] erro HTTP ao acessar CoinGecko: %v", symbol, err)
					continue
				}

				func() {
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						log.Printf("❌ [%s] CoinGecko retornou status %d", symbol, resp.StatusCode)
						return
					}

					var data struct {
						MarketData struct {
							CurrentPrice map[string]float64 `json:"current_price"`
							ATH          map[string]float64 `json:"ath"`
						} `json:"market_data"`
					}

					if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
						log.Printf("❌ [%s] erro ao decodificar resposta: %v", symbol, err)
						return
					}

					current := data.MarketData.CurrentPrice["usd"]
					ath := data.MarketData.ATH["usd"]
					key := strings.ToUpper(symbol)
					last := records[key]

					if current > ath && current > last.AllTimeHigh {
						log.Printf("🚀 [%s] quebrou o recorde histórico oficial! $%.2f > ATH $%.2f", symbol, current, ath)

						msg := GetCryptoPriceMessage(symbol, current, ath)

						alert := fmt.Sprintf("🚨 *NOVO RECORD HISTÓRICO (ATH)*\n\n%s\n\n🕒 ATH superado em %s", msg, time.Now().Format("02/01/2006 15:04"))
						sendAlert(alert)

						records[key] = CryptoRecord{
							AllTimeHigh: current,
							Timestamp:   time.Now(),
						}
						saveRecords(records)
					} else {
						log.Printf("ℹ️ [%s] U$ %.2f — abaixo do ATH oficial U$ %.2f", symbol, current, ath)
					}
				}()
			}
		}
	}()
}

// ResolveAlias converte "BTC" em "bitcoin" usando o mapa PredefinedAliases
func ResolveAlias(symbol string) string {
	alias := strings.ToLower(strings.TrimSpace(symbol))
	if id, ok := PredefinedAliases[alias]; ok {
		return id
	}
	return alias
}

func loadRecords() map[string]CryptoRecord {
	data, err := os.ReadFile(recordFilePath)
	if err != nil {
		log.Println("📂 Nenhum histórico salvo, iniciando novo...")
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

// GetCryptoPriceMessage é um fallback caso não use OpenAI
func GetCryptoPriceMessage(symbol string, current float64, ath float64) string {
	return fmt.Sprintf("*%s*: preço atual U$ %.2f (ATH anterior: U$ %.2f)", strings.ToUpper(symbol), current, ath)
}
