package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// ANSI Color Escape Sequences
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorRed    = "\033[91m"
	ColorGreen  = "\033[92m"
	ColorYellow = "\033[93m"
	ColorBlue   = "\033[94m"
	ColorCyan   = "\033[96m"
	ColorGray   = "\033[90m"
)

// Color variables that can be cleared if output is not a TTY
var (
	colorReset  = ColorReset
	colorBold   = ColorBold
	colorRed    = ColorRed
	colorGreen  = ColorGreen
	colorYellow = ColorYellow
	colorBlue   = ColorBlue
	colorCyan   = ColorCyan
	colorGray   = ColorGray
)

// Supported Fiat currencies
var mapFiat = map[string][2]string{
	"usd": {"USD-BRL", "Dólar Americano"},
	"eur": {"EUR-BRL", "Euro"},
	"gbp": {"GBP-BRL", "Libra Esterlina"},
	"ars": {"ARS-BRL", "Peso Argentino"},
	"jpy": {"JPY-BRL", "Iene Japonês"},
	"cad": {"CAD-BRL", "Dólar Canadense"},
	"aud": {"AUD-BRL", "Dólar Australiano"},
	"chf": {"CHF-BRL", "Franco Suíço"},
	"cny": {"CNY-BRL", "Yuan Chinês"},
}

// Supported Cryptocurrencies
var mapCrypto = map[string][2]string{
	"btc":  {"bitcoin", "Bitcoin"},
	"eth":  {"ethereum", "Ethereum"},
	"sol":  {"solana", "Solana"},
	"ada":  {"cardano", "Cardano"},
	"doge": {"dogecoin", "Dogecoin"},
	"xrp":  {"ripple", "Ripple"},
}

// Synonym mapping to make the user input experience extremely smooth
var synonyms = map[string]string{
	// Fiat synonyms
	"dolar": "usd",
	"dólar": "usd",
	"usd":   "usd",
	"euro":  "eur",
	"eur":   "eur",
	"libra": "gbp",
	"gbp":   "gbp",
	"peso":  "ars",
	"ars":   "ars",
	"iene":  "jpy",
	"jpy":   "jpy",
	"cad":   "cad",
	"aud":   "aud",
	"chf":   "chf",
	"cny":   "cny",

	// Crypto synonyms
	"bitcoin":  "btc",
	"btc":      "btc",
	"ethereum": "eth",
	"eth":      "eth",
	"solana":   "sol",
	"sol":      "sol",
	"cardano":  "ada",
	"ada":      "ada",
	"dogecoin": "doge",
	"doge":     "doge",
	"ripple":   "xrp",
	"xrp":      "xrp",
}

// Structs for AwesomeAPI response
type FiatData struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type AwesomeAPIResponse map[string]FiatData

// Check if stdout is a terminal/TTY
func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// Disable styling when stdout is redirected or piped
func initColors() {
	if !isTTY() {
		colorReset = ""
		colorBold = ""
		colorRed = ""
		colorGreen = ""
		colorYellow = ""
		colorBlue = ""
		colorCyan = ""
		colorGray = ""
	}
}

// Format numbers using standard Brazilian currency symbols and formatting (1.234,56)
func formatCurrencySmart(valFloat float64, currencySymbol string, defaultDecimals int) string {
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

	// Group digits by thousands
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

func formatBRL(valStr string, defaultDecimals int) string {
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return "R$ " + valStr
	}
	return formatCurrencySmart(val, "R$", defaultDecimals)
}

func formatUSD(val float64, defaultDecimals int) string {
	return formatCurrencySmart(val, "US$", defaultDecimals)
}

// Fetch JSON data via HTTP using native libraries with timeouts and custom user-agent
func fetchJSON(url string, target interface{}) error {
	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// CoinGecko and other APIs block requests without a generic browser User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Termcoin/1.0")

	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return fmt.Errorf("falha de conexão. Verifique sua conexão com a internet.")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("recurso não encontrado na API (HTTP 404).")
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro da API (HTTP %d).", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// Format and print a bordered table line in the terminal
func printBorderLine(leftContent, rightContent string, borderLen int, rightColor string) {
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
		coloredRight = fmt.Sprintf("%s%s%s%s", rightColor, colorBold, rightContent, colorReset)
	} else {
		coloredRight = fmt.Sprintf("%s%s%s", colorBold, rightContent, colorReset)
	}

	fmt.Printf("%s│%s%s%s%s%s│%s\n",
		colorBlue, colorReset, innerLabel, coloredRight, strings.Repeat(" ", neededSpaces), colorBlue, colorReset)
}

func displayFiat(symbol, name string, data FiatData) {
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
		varColor = colorGreen
		formattedVarVal = fmt.Sprintf("+%.4f", math.Abs(varVal))
		formattedPctVal = fmt.Sprintf("+%.2f%%", math.Abs(pctVal))
	} else if varVal < 0 {
		varColor = colorRed
		formattedVarVal = fmt.Sprintf("-%.4f", math.Abs(varVal))
		formattedPctVal = fmt.Sprintf("-%.2f%%", math.Abs(pctVal))
	} else {
		varColor = colorGray
		formattedVarVal = "0,0000"
		formattedPctVal = "0,00%"
	}

	// Swapping decimals for PT-BR
	formattedVarVal = strings.Replace(formattedVarVal, ".", ",", -1)
	formattedPctVal = strings.Replace(formattedPctVal, ".", ",", -1)

	formattedBid := formatBRL(bid, 4)
	formattedAsk := formatBRL(ask, 4)
	formattedHigh := formatBRL(high, 4)
	formattedLow := formatBRL(low, 4)

	title := fmt.Sprintf(" %s (%s) ➔ Real Brasileiro (BRL) ", name, strings.ToUpper(symbol))
	titleRuneCount := utf8.RuneCountInString(title)
	borderLen := titleRuneCount + 4
	if borderLen < 55 {
		borderLen = 55
	}

	fmt.Printf("%s┌%s┐%s\n", colorBlue, strings.Repeat("─", borderLen), colorReset)
	
	paddingLeft := (borderLen - titleRuneCount) / 2
	paddingRight := borderLen - titleRuneCount - paddingLeft
	fmt.Printf("%s│%s%s%s%s%s%s│%s\n",
		colorBlue, strings.Repeat(" ", paddingLeft), colorCyan, colorBold, title, colorReset, strings.Repeat(" ", paddingRight), colorBlue)
	
	fmt.Printf("%s├%s┤%s\n", colorBlue, strings.Repeat("─", borderLen), colorReset)

	printBorderLine("Cotação (Compra):", formattedBid, borderLen, "")
	printBorderLine("Cotação (Venda):", formattedAsk, borderLen, "")
	printBorderLine("Máxima do Dia:", formattedHigh, borderLen, "")
	printBorderLine("Mínima do Dia:", formattedLow, borderLen, "")

	varText := fmt.Sprintf("%s (%s)", formattedVarVal, formattedPctVal)
	printBorderLine("Variação do Dia:", varText, borderLen, varColor)

	printBorderLine("Atualizado em:", dateStr, borderLen, "")
	fmt.Printf("%s└%s┘%s\n", colorBlue, strings.Repeat("─", borderLen), colorReset)
}

func displayCrypto(symbol, name string, data map[string]float64) {
	valUSD := data["usd"]
	valBRL := data["brl"]

	formattedUSD := formatUSD(valUSD, 2)
	formattedBRL := formatBRL(strconv.FormatFloat(valBRL, 'f', -1, 64), 2)

	title := fmt.Sprintf(" %s (%s) ➔ Cotação em Tempo Real ", name, strings.ToUpper(symbol))
	titleRuneCount := utf8.RuneCountInString(title)
	borderLen := titleRuneCount + 4
	if borderLen < 55 {
		borderLen = 55
	}

	fmt.Printf("%s┌%s┐%s\n", colorBlue, strings.Repeat("─", borderLen), colorReset)
	
	paddingLeft := (borderLen - titleRuneCount) / 2
	paddingRight := borderLen - titleRuneCount - paddingLeft
	fmt.Printf("%s│%s%s%s%s%s%s│%s\n",
		colorBlue, strings.Repeat(" ", paddingLeft), colorCyan, colorBold, title, colorReset, strings.Repeat(" ", paddingRight), colorBlue)
	
	fmt.Printf("%s├%s┤%s\n", colorBlue, strings.Repeat("─", borderLen), colorReset)

	printBorderLine("Dólar (USD):", formattedUSD, borderLen, "")
	printBorderLine("Real (BRL):", formattedBRL, borderLen, "")
	printBorderLine("Fonte:", "CoinGecko API", borderLen, colorGray)

	fmt.Printf("%s└%s┘%s\n", colorBlue, strings.Repeat("─", borderLen), colorReset)
}

func printHelp() {
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
		colorCyan, colorBold, colorReset,
		colorBold, colorReset,
		colorBold, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorBold, colorReset,
		colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset,
		colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset,
		colorBold, colorReset,
		colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset,
		colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset,
		colorBold, colorReset,
	)
	fmt.Print(helpText)
}

func main() {
	initColors()

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	arg := strings.ToLower(strings.TrimSpace(os.Args[1]))

	if arg == "--help" || arg == "-h" || arg == "help" {
		printHelp()
		os.Exit(0)
	}

	mapped, exists := synonyms[arg]
	if !exists {
		fmt.Fprintf(os.Stderr, "%s%s[Erro] Moeda ou criptomoeda '%s' não é suportada.%s\n", colorRed, colorBold, arg, colorReset)
		fmt.Fprintf(os.Stderr, "Digite %scot --help%s para ver a lista de moedas suportadas.\n", colorGreen, colorReset)
		os.Exit(1)
	}

	if _, ok := mapFiat[mapped]; ok {
		pairInfo := mapFiat[mapped]
		pair := pairInfo[0]
		name := pairInfo[1]
		url := fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%s", pair)

		var response AwesomeAPIResponse
		err := fetchJSON(url, &response)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s%s[Erro] %s%s\n", colorRed, colorBold, err.Error(), colorReset)
			os.Exit(1)
		}

		key := strings.Replace(pair, "-", "", -1)
		fiatData, exists := response[key]
		if !exists {
			fmt.Fprintf(os.Stderr, "%s%s[Erro] Dados não retornados pela API.%s\n", colorRed, colorBold, colorReset)
			os.Exit(1)
		}

		displayFiat(mapped, name, fiatData)
	} else if _, ok := mapCrypto[mapped]; ok {
		cryptoInfo := mapCrypto[mapped]
		cryptoID := cryptoInfo[0]
		name := cryptoInfo[1]
		url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd,brl", cryptoID)

		// Dynamic generic map to avoid complex struct declarations
		var response map[string]map[string]float64
		err := fetchJSON(url, &response)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s%s[Erro] %s%s\n", colorRed, colorBold, err.Error(), colorReset)
			os.Exit(1)
		}

		cryptoData, exists := response[cryptoID]
		if !exists {
			fmt.Fprintf(os.Stderr, "%s%s[Erro] Dados não retornados pela API.%s\n", colorRed, colorBold, colorReset)
			os.Exit(1)
		}

		displayCrypto(mapped, name, cryptoData)
	}
}
