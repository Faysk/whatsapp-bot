package services

import "testing"

func TestFormatNumberBR(t *testing.T) {
	tests := []struct {
		val      float64
		expected string
	}{
		{0, "0,00"},
		{1, "1,00"},
		{1234.56, "1.234,56"},
		{1000000.99, "1.000.000,99"},
	}

	for _, tt := range tests {
		got := formatNumberBR(tt.val)
		if got != tt.expected {
			t.Errorf("formatNumberBR(%v) = %s, want %s", tt.val, got, tt.expected)
		}
	}
}

func TestFormatNumberUS(t *testing.T) {
	tests := []struct {
		val      float64
		expected string
	}{
		{0, "0.00"},
		{1, "1.00"},
		{1234.56, "1,234.56"},
		{1000000.99, "1,000,000.99"},
	}

	for _, tt := range tests {
		got := formatNumberUS(tt.val)
		if got != tt.expected {
			t.Errorf("formatNumberUS(%v) = %s, want %s", tt.val, got, tt.expected)
		}
	}
}
