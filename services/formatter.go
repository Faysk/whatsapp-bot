// File: services/formatter.go
package services

import (
	"fmt"
	"strconv"
)

// formatNumberBR formata número no estilo brasileiro (1.000.000,00)
func formatNumberBR(val float64) string {
	intPart := int64(val)
	decimal := int64((val - float64(intPart)) * 100)
	intStr := formatWithSeparator(intPart, ".")
	return fmt.Sprintf("%s,%02d", intStr, decimal)
}

// formatNumberUS formata número no estilo americano (1,000,000.00)
func formatNumberUS(val float64) string {
	intPart := int64(val)
	decimal := int64((val - float64(intPart)) * 100)
	intStr := formatWithSeparator(intPart, ",")
	return fmt.Sprintf("%s.%02d", intStr, decimal)
}

// formatWithSeparator aplica separadores de milhar reversamente
func formatWithSeparator(n int64, sep string) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}

	var result []byte
	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if count != 0 && count%3 == 0 {
			result = append(result, sep[0])
		}
		result = append(result, s[i])
		count++
	}

	// inverte resultado
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
