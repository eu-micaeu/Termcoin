package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"cot/internal/currency"
)

// PrintBorderLine prints a single line with visual alignment using rune count
func PrintBorderLine(leftContent, rightContent string, borderLen int, rightColor string) {
	innerLabel := "  " + leftContent
	labelRuneCount := utf8.RuneCountInString(innerLabel)
	if labelRuneCount < 24 {
		innerLabel = innerLabel + strings.Repeat(" ", 24-labelRuneCount)
		labelRuneCount = 24
	}

	neededSpaces := borderLen - labelRuneCount - utf8.RuneCountInString(rightContent)
	if neededSpaces < 0 {
		neededSpaces = 0
	}

	var coloredRight string
	if rightColor != "" {
		coloredRight = fmt.Sprintf("%s%s%s%s", rightColor, Bold, rightContent, Reset)
	} else {
		coloredRight = fmt.Sprintf("%s%s%s", Bold, rightContent, Reset)
	}

	fmt.Printf("%s│%s%s%s%s%s│%s\n",
		Blue, Reset, innerLabel, coloredRight, strings.Repeat(" ", neededSpaces), Blue, Reset)
}

func DisplayFiat(symbol, name string, data currency.FiatData) {
	bid := data.Bid
	ask := data.Ask
	high := data.High
	low := data.Low
	pctChange := data.PctChange
	varBid := data.VarBid
	dateStr := data.CreateDate

	varVal, _ := strconv.ParseFloat(varBid, 64)
	pctVal, _ := strconv.ParseFloat(pctChange, 64)

	var varColor string
	var formattedVarVal string
	var formattedPctVal string

	if varVal > 0 {
		varColor = Green
		formattedVarVal = fmt.Sprintf("+%.4f", math.Abs(varVal))
		formattedPctVal = fmt.Sprintf("+%.2f%%", math.Abs(pctVal))
	} else if varVal < 0 {
		varColor = Red
		formattedVarVal = fmt.Sprintf("-%.4f", math.Abs(varVal))
		formattedPctVal = fmt.Sprintf("-%.2f%%", math.Abs(pctVal))
	} else {
		varColor = Gray
		formattedVarVal = "0,0000"
		formattedPctVal = "0,00%"
	}

	formattedVarVal = strings.Replace(formattedVarVal, ".", ",", -1)
	formattedPctVal = strings.Replace(formattedPctVal, ".", ",", -1)

	formattedBid := currency.FormatBRL(bid, 4)
	formattedAsk := currency.FormatBRL(ask, 4)
	formattedHigh := currency.FormatBRL(high, 4)
	formattedLow := currency.FormatBRL(low, 4)

	title := fmt.Sprintf(" %s (%s) ➔ Real Brasileiro (BRL) ", name, strings.ToUpper(symbol))
	titleRuneCount := utf8.RuneCountInString(title)
	borderLen := titleRuneCount + 4
	if borderLen < 55 {
		borderLen = 55
	}

	fmt.Printf("%s┌%s┐%s\n", Blue, strings.Repeat("─", borderLen), Reset)
	
	paddingLeft := (borderLen - titleRuneCount) / 2
	paddingRight := borderLen - titleRuneCount - paddingLeft
	fmt.Printf("%s│%s%s%s%s%s%s│%s\n",
		Blue, strings.Repeat(" ", paddingLeft), Cyan, Bold, title, Reset, strings.Repeat(" ", paddingRight), Blue)
	
	fmt.Printf("%s├%s┤%s\n", Blue, strings.Repeat("─", borderLen), Reset)

	PrintBorderLine("Cotação (Compra):", formattedBid, borderLen, "")
	PrintBorderLine("Cotação (Venda):", formattedAsk, borderLen, "")
	PrintBorderLine("Máxima do Dia:", formattedHigh, borderLen, "")
	PrintBorderLine("Mínima do Dia:", formattedLow, borderLen, "")

	varText := fmt.Sprintf("%s (%s)", formattedVarVal, formattedPctVal)
	PrintBorderLine("Variação do Dia:", varText, borderLen, varColor)

	PrintBorderLine("Atualizado em:", dateStr, borderLen, "")
	fmt.Printf("%s└%s┘%s\n", Blue, strings.Repeat("─", borderLen), Reset)
}

func DisplayCrypto(symbol, name string, data map[string]float64) {
	valUSD := data["usd"]
	valBRL := data["brl"]

	formattedUSD := currency.FormatUSD(valUSD, 2)
	formattedBRL := currency.FormatBRL(strconv.FormatFloat(valBRL, 'f', -1, 64), 2)

	title := fmt.Sprintf(" %s (%s) ➔ Cotação em Tempo Real ", name, strings.ToUpper(symbol))
	titleRuneCount := utf8.RuneCountInString(title)
	borderLen := titleRuneCount + 4
	if borderLen < 55 {
		borderLen = 55
	}

	fmt.Printf("%s┌%s┐%s\n", Blue, strings.Repeat("─", borderLen), Reset)
	
	paddingLeft := (borderLen - titleRuneCount) / 2
	paddingRight := borderLen - titleRuneCount - paddingLeft
	fmt.Printf("%s│%s%s%s%s%s%s│%s\n",
		Blue, strings.Repeat(" ", paddingLeft), Cyan, Bold, title, Reset, strings.Repeat(" ", paddingRight), Blue)
	
	fmt.Printf("%s├%s┤%s\n", Blue, strings.Repeat("─", borderLen), Reset)

	PrintBorderLine("Dólar (USD):", formattedUSD, borderLen, "")
	PrintBorderLine("Real (BRL):", formattedBRL, borderLen, "")
	PrintBorderLine("Fonte:", "CoinGecko API", borderLen, Gray)

	fmt.Printf("%s└%s┘%s\n", Blue, strings.Repeat("─", borderLen), Reset)
}

func PrintHelp() {
	helpText := fmt.Sprintf(`
%s%sCot - Monitor de Cotações em Tempo Real (CLI)%s

%sUso:%s
  cot <moeda>

%sExemplos:%s
  %scot usd%s     - Exibe cotação do Dólar Americano em Real (BRL)
  %scot eur%s     - Exibe cotação do Euro em Real (BRL)
  %scot btc%s     - Exibe cotação do Bitcoin em Dólar (USD) e Real (BRL)
  %scot eth%s     - Exibe cotação do Ethereum em Dólar (USD) e Real (BRL)

%sMoedas Tradicionais (Fiat) Suportadas:%s
  %susd%s (Dólar), %seur%s (Euro), %sgbp%s (Libra), %sars%s (Peso Arg), %sjpy%s (Iene), 
  %scad%s (Dólar Can), %saud%s (Dólar Aus), %schf%s (Franco Suíço), %scny%s (Yuan)

%sCriptomoedas Suportadas:%s
  %sbtc%s (Bitcoin), %seth%s (Ethereum), %ssol%s (Solana), 
  %sada%s (Cardano), %sdoge%s (Dogecoin), %sxrp%s (Ripple)

%sAjuda:%s
  cot --help  - Mostra esta tela de ajuda
`,
		Cyan, Bold, Reset,
		Bold, Reset,
		Bold, Reset,
		Green, Reset,
		Green, Reset,
		Green, Reset,
		Green, Reset,
		Bold, Reset,
		Yellow, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset,
		Yellow, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset,
		Bold, Reset,
		Yellow, Reset, Yellow, Reset, Yellow, Reset,
		Yellow, Reset, Yellow, Reset, Yellow, Reset,
		Bold, Reset,
	)
	fmt.Print(helpText)
}
