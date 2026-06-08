package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// FetchJSON handles basic client configuration, User-Agent, and decoding JSON
func FetchJSON(url string, target interface{}) error {
	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
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

// BinanceResponse holds the symbol and price from Binance API ticker
type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// FetchBinancePrice queries both USDT and BRL prices for a crypto asset
func FetchBinancePrice(symbolUSDT, symbolBRL string) (float64, float64, error) {
	var resUSDT BinanceResponse
	err := FetchJSON(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbolUSDT), &resUSDT)
	if err != nil {
		return 0, 0, err
	}

	var resBRL BinanceResponse
	err = FetchJSON(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbolBRL), &resBRL)
	if err != nil {
		return 0, 0, err
	}

	pUSD, _ := strconv.ParseFloat(resUSDT.Price, 64)
	pBRL, _ := strconv.ParseFloat(resBRL.Price, 64)

	return pUSD, pBRL, nil
}
