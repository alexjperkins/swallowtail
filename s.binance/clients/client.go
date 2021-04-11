package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	// Endpoints
	binanceAPIURL        = "https://api.binance.com"
	binanceBackupAPIURLs = []string{
		"https://api1.binance.com",
		"https://api2.binance.com",
		"https://api3.binance.com",
	}
	testBinanceAPIEndpoint = ""

	binanceSpotEndpoint     = "/api/v3/order"
	testBinanceSpotEndpoint = "/api/v3/order/test"

	binanceLimitOrderID  = "LIMIT"
	binanceMarketOrderID = "MARKET"
	// TODO check
	binanceStopLossID = "STOP-LOSS"
)

func New(timeout time.Duration) *BinanceClient {
	return &BinanceClient{
		c: &http.Client{
			Timeout: timeout,
		},
	}
}

// NOTE: Auth passed with X-MBX-APIKEY
type BinanceClient struct {
	c *http.Client
}

type spotOrder struct {
	symbol      string
	side        string
	tradeType   string `json:"type"`
	timeInForce string
	quantity    float32
	price       float32
	stopPrice   float32
	timestamp   time.Time
}

func newSpotOrder(symbol, side, tradeType, timeInForce string, quantity, price, stopPrice float32, timestamp time.Time) *spotOrder {
	return &spotOrder{
		symbol:      symbol,
		side:        side,
		tradeType:   tradeType,
		timeInForce: timeInForce,
		quantity:    quantity,
		price:       price,
		stopPrice:   stopPrice,
		timestamp:   timestamp,
	}
}

func (b *BinanceClient) spotOrder(symbol, side, tradeType, timeInForce string, quantity, entry, sl float32, tp ...float32) error {
	// validation of params
	spotOrder := newSpotOrder(symbol, side, tradeType, timeInForce, quantity, entry, sl, time.Now())
	reqBody, err := json.Marshal(spotOrder)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s%s", binanceAPIURL, binanceSpotEndpoint)
	rsp, err := b.c.Post(endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	return nil
}

func (b *BinanceClient) futuresOrder(symbol, side, typeOfTrade string, riskPercentage, entry, sl float32, tp ...float32) error {
	return nil
}

func (b *BinanceClient) SpotLimitOrder(symbol, side, timeInForce string, riskPercentage, entry, sl float32, tp ...float32) error {
	// calculate risk percentage
	var (
		quantity float32
	)
	err := b.spotOrder(symbol, side, binanceLimitOrderID, timeInForce, quantity, entry, sl)
	if err != nil {
		return err
	}
	// Take Profits
	return nil
}

func (b *BinanceClient) SpotDCAOrder(ticker, side string, riskPercentage, upper, lower, sl float32, tp ...float32) error {
	return nil
}

func (b *BinanceClient) SpotMarketOrder(ticker, side string, riskPercentage, sl float32, checkPriceDistanceBeforeEntering bool, tp ...float32) error {
	return nil
}

func (b *BinanceClient) SpotGetPrice(ticker string) (float32, error) {
	return 0.0, nil
}

func (b *BinanceClient) FuturesLimitOrder(ticker, side string, riskPercentage, entry, sl float32, tp ...float32) error {
	return nil
}

func (b *BinanceClient) FuturesDCAOrder(ticker, side string, riskPercentage, upper, lower, sl float32, tp ...float32) error {
	return nil
}

func (b *BinanceClient) FuturesMarketOrder(ticker, side string, riskPercentage, sl float32, checkPriceDistanceBeforeEntering bool, tp ...float32) error {
	return nil
}

func (b *BinanceClient) FuturesGetPrice(ticker string) (float32, error) {
	return 0.0, nil
}

func (b *BinanceClient) GetAccountSize() (float32, error) {
	return 0.0, nil
}
