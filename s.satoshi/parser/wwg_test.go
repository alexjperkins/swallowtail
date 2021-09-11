package parser

import (
	"context"
	tradeengineproto "swallowtail/s.trade-engine/proto"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWWGParser(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		username      string
		truthValue    float64
		expectedTrade *tradeengineproto.Trade
		withErr       bool
	}{
		{
			name: "rego_full_trade_wwg_with_three_tp",
			content: `Hey guys I'm LONG BTC here.

			ENTRY: 50000
			STOP: 49000

			TP1: 52000
			TP2: 54000
			TP3: 58000

			This should give us an 4.5RR and 15.7% increase
			`,
			username:   "rego",
			truthValue: 50000,
			expectedTrade: &tradeengineproto.Trade{
				ActorId:   "REGO",
				ActorType: tradeengineproto.ACTOR_TYPE_EXTERNAL.String(),
				OrderType: tradeengineproto.ORDER_TYPE_MARKET,
				Asset:     "BTC",
				Pair:      tradeengineproto.TRADE_PAIR_USDT,
				TradeSide: tradeengineproto.TRADE_SIDE_LONG,
				TradeType: tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
				Entry:     50000,
				StopLoss:  49000,
				TakeProfits: []float32{
					52000, 54000, 58000,
				},
			},
		},
		{
			name:       "bluntz_example_second_entry",
			username:   "bluntz",
			truthValue: 170,
			content: `Going to enter that second sol entry here as i think it just got frontrun by 0.3%

			entry 2: now 165
			stop 135.61
			target 259.7

			57% 3.25RR
			`,
			expectedTrade: &tradeengineproto.Trade{
				ActorId:   "BLUNTZ",
				ActorType: tradeengineproto.ACTOR_TYPE_EXTERNAL.String(),
				OrderType: tradeengineproto.ORDER_TYPE_MARKET,
				Asset:     "SOL",
				Pair:      tradeengineproto.TRADE_PAIR_USDT,
				TradeSide: tradeengineproto.TRADE_SIDE_LONG,
				TradeType: tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
				Entry:     165,
				StopLoss:  135.61,
				TakeProfits: []float32{
					259.7,
				},
			},
		},
		{
			name:       "astekz_example_1_aave_no_take_profit",
			username:   "astekz",
			truthValue: 344,
			content: `
			aave 
			spot or low lev long 343
			stop 323
			@​everyone
			[Attachments]
			https://cdn.discordapp.com/attachments/869596440777883749/885529381479518219/unknown.png
			`,
			expectedTrade: &tradeengineproto.Trade{
				ActorId:     "ASTEKZ",
				ActorType:   tradeengineproto.ACTOR_TYPE_EXTERNAL.String(),
				OrderType:   tradeengineproto.ORDER_TYPE_MARKET,
				Asset:       "AAVE",
				Pair:        tradeengineproto.TRADE_PAIR_USDT,
				TradeSide:   tradeengineproto.TRADE_SIDE_LONG,
				TradeType:   tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
				Entry:       343,
				StopLoss:    323,
				TakeProfits: []float32{},
			},
		},
		{
			name:       "eli_example_1_limit_srm",
			username:   "eli",
			truthValue: 10.9,
			content:    `SRM LIMIT LONG 9.80 stop 8.90 tp 13 18 @​everyone`,
			expectedTrade: &tradeengineproto.Trade{
				ActorId:     "ELI",
				ActorType:   tradeengineproto.ACTOR_TYPE_EXTERNAL.String(),
				OrderType:   tradeengineproto.ORDER_TYPE_LIMIT,
				Asset:       "SRM",
				Pair:        tradeengineproto.TRADE_PAIR_USDT,
				TradeSide:   tradeengineproto.TRADE_SIDE_LONG,
				TradeType:   tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
				Entry:       9.80,
				StopLoss:    8.90,
				TakeProfits: []float32{13, 18},
			},
		},
		{
			name:       "cryptogodjohnny_example_1_market_buy_srm",
			username:   "cryptogodjohnny",
			truthValue: 0.041,
			content: `
			RSR $0.0402

			SL $0.0374
			`,
			expectedTrade: &tradeengineproto.Trade{
				ActorId:     "CRYPTOGODJOHNNY",
				ActorType:   tradeengineproto.ACTOR_TYPE_EXTERNAL.String(),
				OrderType:   tradeengineproto.ORDER_TYPE_MARKET,
				Asset:       "RSR",
				Pair:        tradeengineproto.TRADE_PAIR_USDT,
				TradeSide:   tradeengineproto.TRADE_SIDE_LONG,
				TradeType:   tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
				Entry:       0.0402,
				StopLoss:    0.0374,
				TakeProfits: []float32{},
			},
		},
		{
			name:       "cryptogodjohnny_example_2_market_btc_short",
			username:   "cryptogodjohnny",
			truthValue: 46500,
			content: `
			Btc short $46650

			SL 47801 

			Tp 45800 44540 43680 42112
			@​Scalps High risk
			`,
			expectedTrade: &tradeengineproto.Trade{
				ActorId:   "CRYPTOGODJOHNNY",
				ActorType: tradeengineproto.ACTOR_TYPE_EXTERNAL.String(),
				OrderType: tradeengineproto.ORDER_TYPE_MARKET,
				Asset:     "BTC",
				Pair:      tradeengineproto.TRADE_PAIR_USDT,
				TradeSide: tradeengineproto.TRADE_SIDE_SHORT,
				TradeType: tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
				Entry:     46650,
				StopLoss:  47801,
				TakeProfits: []float32{
					45800,
					44540,
					43680,
					42112,
				},
			},
		},
		{
			name:          "ticker_but_no_valid_information_example_ftt",
			content:       `if i ever get ftt at 50 again im gonna put entire portfolio there like jeliaz said`,
			expectedTrade: nil,
			withErr:       true,
		},
	}

	originalBinanceAssetPairs := binanceAssetPairs
	binanceAssetPairs = map[string]bool{
		"btc":  true,
		"sol":  true,
		"aave": true,
		"srm":  true,
		"rsr":  true,
		"ftt":  true,
	}

	originalFetcher := fetchLatestPrice
	t.Cleanup(func() {
		binanceAssetPairs = originalBinanceAssetPairs
		fetchLatestPrice = originalFetcher
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetchLatestPrice = func(_ context.Context, _, _ string) (float64, error) {
				return tt.truthValue, nil
			}

			trade, err := (&WWGParser{}).Parse(context.Background(), tt.content, &discordgo.MessageCreate{
				Message: &discordgo.Message{
					Author: &discordgo.User{
						Username: tt.username,
					},
				},
			})

			switch {
			case !tt.withErr:
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTrade, trade)
			default:
				require.Error(t, err)
			}
		})
	}
}
