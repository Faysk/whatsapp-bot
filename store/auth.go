package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/faysk/whatsapp-bot/config"
)

const authorizedPath = "authorized.json"

var phoneRegex = regexp.MustCompile(`^55\d{10,11}$`)

// LoadAuthorizedNumbers carrega e une números fixos e dinâmicos (com validação)
func LoadAuthorizedNumbers() []string {
	data, err := os.ReadFile(authorizedPath)
	if err != nil {
		log.Printf("⚠️ Arquivo %s não encontrado. Criando novo com lista vazia.", authorizedPath)
		_ = SaveAuthorizedNumbers([]string{})
		return mergeWithFixed([]string{})
	}

	var dynamic []string
	if err := json.Unmarshal(data, &dynamic); err != nil {
		log.Printf("❌ Erro ao decodificar %s: %v. Substituindo por lista vazia.", authorizedPath, err)
		_ = SaveAuthorizedNumbers([]string{})
		return mergeWithFixed([]string{})
	}

	return mergeWithFixed(dynamic)
}

// SaveAuthorizedNumbers salva apenas os números mutáveis (exclui fixos)
func SaveAuthorizedNumbers(all []string) error {
	mutables := filterMutable(all)
	mutables = sanitize(mutables)

	data, err := json.MarshalIndent(mutables, "", "  ")
	if err != nil {
		return fmt.Errorf("❌ Erro ao gerar JSON: %w", err)
	}

	if err := os.WriteFile(authorizedPath, data, 0644); err != nil {
		return fmt.Errorf("❌ Erro ao salvar %s: %w", authorizedPath, err)
	}

	log.Printf("✅ Lista de autorizados salva com %d número(s) mutáveis.", len(mutables))
	return nil
}

// AddAuthorized adiciona um número à lista, se válido, não fixo e não duplicado
func AddAuthorized(num string) error {
	num = strings.TrimSpace(num)

	switch {
	case num == "":
		return fmt.Errorf("⚠️ Número vazio ignorado.")
	case IsFixed(num):
		return fmt.Errorf("⚠️ Número %s é fixo, não pode ser adicionado via comando.", num)
	case !isValidPhone(num):
		return fmt.Errorf("⚠️ Número inválido: %s", num)
	}

	list := LoadAuthorizedNumbers()
	if contains(list, num) {
		log.Printf("ℹ️ Número %s já estava autorizado.", num)
		return nil
	}

	list = append(list, num)
	return SaveAuthorizedNumbers(list)
}

// RemoveAuthorized remove um número, se não for fixo nem o próprio solicitante
func RemoveAuthorized(requester, target string) error {
	switch {
	case requester == target:
		return fmt.Errorf("⚠️ %s tentou se auto-remover. Operação bloqueada.", requester)
	case IsFixed(target):
		return fmt.Errorf("⚠️ Tentativa de remover número fixo %s foi bloqueada.", target)
	}

	list := LoadAuthorizedNumbers()
	var updated []string
	for _, n := range list {
		if n != target {
			updated = append(updated, n)
		}
	}

	if len(updated) == len(list) {
		log.Printf("ℹ️ Número %s não estava na lista. Nenhuma alteração feita.", target)
	} else {
		log.Printf("🗑️ Número %s removido com sucesso.", target)
	}

	return SaveAuthorizedNumbers(updated)
}

// IsFixed verifica se o número é fixo e não pode ser removido
func IsFixed(num string) bool {
	return contains(config.AppConfig.FixedAuthorizedEnv, num)
}

//
// === 🧠 Utilitários Internos ===
//

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
		if !IsFixed(n) {
			result = append(result, n)
		}
	}
	return result
}

func sanitize(list []string) []string {
	unique := map[string]struct{}{}
	for _, n := range list {
		n = strings.TrimSpace(n)
		if isValidPhone(n) {
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
