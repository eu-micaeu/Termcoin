package currency

import "strings"

// Supported Fiat currencies
var MapFiat = map[string][2]string{
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
var MapCrypto = map[string][2]string{
	"btc":  {"bitcoin", "Bitcoin"},
	"eth":  {"ethereum", "Ethereum"},
	"sol":  {"solana", "Solana"},
	"ada":  {"cardano", "Cardano"},
	"doge": {"dogecoin", "Dogecoin"},
	"xrp":  {"ripple", "Ripple"},
}

// Synonym mapping to make the user input experience extremely smooth
var Synonyms = map[string]string{
	"dolar":    "usd",
	"dólar":    "usd",
	"usd":      "usd",
	"euro":     "eur",
	"eur":      "eur",
	"libra":    "gbp",
	"gbp":      "gbp",
	"peso":     "ars",
	"ars":      "ars",
	"iene":     "jpy",
	"jpy":      "jpy",
	"cad":      "cad",
	"aud":      "aud",
	"chf":      "chf",
	"cny":      "cny",
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

// Resolve matches a user input to a canonical currency symbol (e.g. "bitcoin" -> "btc")
func Resolve(input string) (string, bool) {
	clean := strings.ToLower(strings.TrimSpace(input))
	canonical, exists := Synonyms[clean]
	return canonical, exists
}
