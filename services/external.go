package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// init carrega os dados das moedas e configura os aliases dinâmicos
func init() {
	// Carrega aliases fixos do arquivo crypto_aliases.go
	for alias, id := range PredefinedAliases {
		cryptoAliases[alias] = id
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.coingecko.com/api/v3/coins/list")
	if err != nil {
		fmt.Println("⚠️ Erro ao consultar CoinGecko:", err)
		return
	}
	defer resp.Body.Close()

	var coins []struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	}

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

		addAliasIfAbsent(id, id)
		addAliasIfAbsent(symbol, id)
		addAliasIfAbsent(name, id)
		addAliasIfAbsent(strings.ReplaceAll(name, " ", ""), id)
	}
}

// addAliasIfAbsent adiciona o alias somente se ainda não existir
func addAliasIfAbsent(alias, id string) {
	if _, exists := cryptoAliases[alias]; !exists {
		cryptoAliases[alias] = id
	}
}
