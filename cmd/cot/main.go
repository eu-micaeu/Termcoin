package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"cot/internal/api"
	"cot/internal/currency"
	"cot/internal/ui"
)

func main() {
	ui.InitColors()

	if len(os.Args) < 2 {
		ui.PrintHelp()
		os.Exit(0)
	}

	arg := strings.ToLower(strings.TrimSpace(os.Args[1]))

	if arg == "--help" || arg == "-h" || arg == "help" {
		ui.PrintHelp()
		os.Exit(0)
	}

	mapped, exists := currency.Resolve(arg)
	if !exists {
		fmt.Fprintf(os.Stderr, "%s%s[Erro] Moeda ou criptomoeda '%s' não é suportada.%s\n", ui.Red, ui.Bold, arg, ui.Reset)
		fmt.Fprintf(os.Stderr, "Digite %scot --help%s para ver a lista de moedas suportadas.\n", ui.Green, ui.Reset)
		os.Exit(1)
	}

	runRealTimeLoop(mapped)
}

func runRealTimeLoop(mapped string) {
	// Signal interception for restoring the cursor on Ctrl+C / SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Print("\033[?25h") // restore cursor
		fmt.Print("\nSaindo...\n")
		os.Exit(0)
	}()

	// Hide terminal cursor for a cleaner UI
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	var prices []float64

	isFiat := false
	var pair, name, symbolUSDT, symbolBRL string
	var currencySymbol string

	if val, ok := currency.MapFiat[mapped]; ok {
		isFiat = true
		pair = val[0]
		name = val[1]
		currencySymbol = "R$"
	} else if val, ok := currency.MapCrypto[mapped]; ok {
		name = val[1]
		symbolUSDT = strings.ToUpper(mapped) + "USDT"
		symbolBRL = strings.ToUpper(mapped) + "BRL"
		currencySymbol = "US$"
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Track the number of lines printed to overwrite them precisely on the next tick
	var linesPrinted int

	// Initial immediate fetch so user doesn't wait 1 second for the first frame
	linesPrinted = fetchAndDraw(isFiat, pair, symbolUSDT, symbolBRL, mapped, name, currencySymbol, &prices, linesPrinted)

	for {
		select {
		case <-ticker.C:
			linesPrinted = fetchAndDraw(isFiat, pair, symbolUSDT, symbolBRL, mapped, name, currencySymbol, &prices, linesPrinted)
		}
	}
}

func fetchAndDraw(isFiat bool, pair, symbolUSDT, symbolBRL, mapped, name, currencySymbol string, prices *[]float64, prevLines int) int {
	var currentPrice float64
	var currentBRL float64
	var err error

	if isFiat {
		url := fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%s", pair)
		var response currency.AwesomeAPIResponse
		err = api.FetchJSON(url, &response)
		if err == nil {
			key := strings.Replace(pair, "-", "", -1)
			if data, ok := response[key]; ok {
				currentPrice, _ = strconv.ParseFloat(data.Bid, 64)
			} else {
				err = fmt.Errorf("dados indisponíveis")
			}
		}
	} else {
		currentPrice, currentBRL, err = api.FetchBinancePrice(symbolUSDT, symbolBRL)
	}

	// Move cursor up to overwrite previous draw
	if prevLines > 0 {
		fmt.Printf("\033[%dA", prevLines)
	}

	if err != nil {
		fmt.Printf("\033[K%s%s[Erro ao atualizar]: %s%s\n", ui.Red, ui.Bold, err.Error(), ui.Reset)
		return 1
	}

	*prices = append(*prices, currentPrice)
	// Limit to 45 points to fit typical terminal width nicely (45 + 15 label width = 60 chars)
	if len(*prices) > 45 {
		*prices = (*prices)[1:]
	}

	if len(*prices) > 0 {
		return ui.DrawChart(*prices, mapped, name, currencySymbol, currentBRL)
	}
	return 0
}
