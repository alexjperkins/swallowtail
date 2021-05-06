package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	BinanceHTTPClientID = "binance-client"

	binanceHostname           = "https://fapi.binance.com"
	binancePingEndpoint       = "/fapi/v1/ping"
	binanceKlineStickEndpoint = "/fapi/v1/klines"
	binanceMarkPriceEndpoint  = "/fapi/v1/premiumIndex"
)

func NewBinanceClient() *BinanceClient {
	return &BinanceClient{
		c: &http.Client{},
	}
}

type BinanceClient struct {
	c *http.Client
}

type BinanceMarkPriceRsp struct {
	Symbol          string `json:"symbol"`
	MarkPrice       string `json:"markPrice"`
	IndexPrice      string `json:"indexPrice"`
	LastFundingRate string `json:"lastFundingRate"`
	NextFundingTime int    `json:"nextFundingTime"`
	InterestRate    string `json:"interestRate"`
	Time            int    `json:"time"`
}

func (bc *BinanceClient) GetPrice(ctx context.Context, symbol string) (*BinanceMarkPriceRsp, error) {
	rsp, err := bc.c.Get(fmt.Sprintf("%s%s?symbol=%s", binanceHostname, binanceMarkPriceEndpoint, symbol))
	if err != nil {
		return nil, nil
	}
	defer rsp.Body.Close()

	var binanceMarkPriceRsp *BinanceMarkPriceRsp
	err = json.NewDecoder(rsp.Body).Decode(&binanceMarkPriceRsp)
	if err != nil {
		return nil, err
	}
	return binanceMarkPriceRsp, nil
}

func (bc *BinanceClient) Ping() error {
	_, err := bc.c.Get(fmt.Sprintf("%s%s", binanceHostname, binancePingEndpoint))
	return err
}
