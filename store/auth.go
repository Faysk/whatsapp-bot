package store

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/faysk/whatsapp-bot/config"
)

const authorizedPath = "authorized.json"
var phoneRegex = regexp.MustCompile(`^55\d{10,11}$`)

// LoadAuthorizedNumbers carrega os números do JSON e mescla com os fixos do .env
func LoadAuthorizedNumbers() []string {
	data, err := os.ReadFile(authorizedPath)
	if err != nil {
		log.Printf("⚠️ Arquivo %s não encontrado. Criando novo com lista vazia.", authorizedPath)
		_ = SaveAuthorizedNumbers([]string{})
		return mergeWithFixed([]string{})
	}

	var dynamic []string
	if err := json.Unmarshal(data, &dynamic); err != nil {
		log.Printf("❌ Erro ao ler %s: %v", authorizedPath, err)
		return mergeWithFixed([]string{})
	}

	return mergeWithFixed(dynamic)
}

// SaveAuthorizedNumbers salva apenas os números mutáveis (excluindo os fixos)
func SaveAuthorizedNumbers(all []string) error {
	mutables := filterMutable(all)
	mutables = sanitize(mutables)

	data, err := json.MarshalIndent(mutables, "", "  ")
	if err != nil {
		log.Printf("❌ Erro ao gerar JSON: %v", err)
		return err
	}

	if err := os.WriteFile(authorizedPath, data, 0644); err != nil {
		log.Printf("❌ Erro ao salvar %s: %v", authorizedPath, err)
		return err
	}

	log.Printf("✅ Lista salva com sucesso (%d número(s) mutáveis).", len(mutables))
	return nil
}

// AddAuthorized adiciona um novo número à lista (se não for fixo ou duplicado)
func AddAuthorized(num string) error {
	num = strings.TrimSpace(num)
	if num == "" || IsFixed(num) {
		log.Printf("⚠️ Número %s ignorado (vazio ou fixo).", num)
		return nil
	}

	if !isValidPhone(num) {
		log.Printf("⚠️ Número inválido ignorado: %s", num)
		return nil
	}

	list := LoadAuthorizedNumbers()
	if contains(list, num) {
		log.Printf("ℹ️ Número %s já está autorizado. Nenhuma alteração.", num)
		return nil
	}

	list = append(list, num)
	return SaveAuthorizedNumbers(list)
}

// RemoveAuthorized remove um número, se não for fixo e não for o próprio solicitante
func RemoveAuthorized(requester, target string) error {
	if requester == target {
		log.Printf("⚠️ %s tentou se remover da lista — operação ignorada.", requester)
		return nil
	}
	if IsFixed(target) {
		log.Printf("⚠️ Tentativa de remover número fixo %s — bloqueado.", target)
		return nil
	}

	list := LoadAuthorizedNumbers()
	var updated []string
	for _, n := range list {
		if n != target {
			updated = append(updated, n)
		}
	}

	if len(updated) == len(list) {
		log.Printf("ℹ️ Número %s não estava na lista. Nenhuma alteração.", target)
	} else {
		log.Printf("🗑️ Número %s removido da lista de autorizados.", target)
	}

	return SaveAuthorizedNumbers(updated)
}

// IsFixed verifica se o número está entre os fixos do .env
func IsFixed(num string) bool {
	return contains(config.AppConfig.FixedAuthorizedEnv, num)
}

// === Utilitários ===

func mergeWithFixed(mutables []string) []string {
	all := append([]string{}, config.AppConfig.FixedAuthorizedEnv...)
	for _, n := range mutables {
		if !contains(all, n) {
			all = append(all, n)
		}
	}
	return sanitize(all)
}

func filterMutable(all []string) []string {
	var result []string
	for _, n := range all {
		if !contains(config.AppConfig.FixedAuthorizedEnv, n) {
			result = append(result, n)
		}
	}
	return result
}

func sanitize(list []string) []string {
	unique := map[string]struct{}{}
	for _, n := range list {
		n = strings.TrimSpace(n)
		if n != "" {
			unique[n] = struct{}{}
		}
	}
	var result []string
	for n := range unique {
		result = append(result, n)
	}
	sort.Strings(result)
	return result
}

func contains(list []string, value string) bool {
	for _, n := range list {
		if n == value {
			return true
		}
	}
	return false
}

func isValidPhone(num string) bool {
	return phoneRegex.MatchString(num)
}
