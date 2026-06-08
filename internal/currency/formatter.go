package currency

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// FormatCurrencySmart formats a number with PT-BR commas and periods, dynamically choosing precision
func FormatCurrencySmart(valFloat float64, currencySymbol string, defaultDecimals int) string {
	var decimals int
	absVal := math.Abs(valFloat)

	if valFloat == 0.0 {
		decimals = 2
	} else if absVal < 0.0001 {
		decimals = 8
	} else if absVal < 0.01 {
		decimals = 6
	} else if absVal < 1.0 {
		decimals = 4
	} else {
		decimals = defaultDecimals
	}

	isNegative := valFloat < 0
	raw := fmt.Sprintf("%.*f", decimals, absVal)

	parts := strings.Split(raw, ".")
	integerPart := parts[0]
	var decimalPart string
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	var result []string
	length := len(integerPart)
	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{integerPart[start:i]}, result...)
	}
	formattedInt := strings.Join(result, ".")

	sign := ""
	if isNegative {
		sign = "-"
	}

	if decimalPart != "" {
		return fmt.Sprintf("%s %s%s,%s", currencySymbol, sign, formattedInt, decimalPart)
	}
	return fmt.Sprintf("%s %s%s", currencySymbol, sign, formattedInt)
}

func FormatBRL(valStr string, defaultDecimals int) string {
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return "R$ " + valStr
	}
	return FormatCurrencySmart(val, "R$", defaultDecimals)
}

func FormatUSD(val float64, defaultDecimals int) string {
	return FormatCurrencySmart(val, "US$", defaultDecimals)
}
