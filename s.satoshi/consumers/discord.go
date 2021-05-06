package consumers

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"swallowtail/libraries/util"
	coingecko "swallowtail/s.coingecko/clients"
	"swallowtail/s.satoshi/clients"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/monzo/slog"
)

var (
	PostToDiscordConsumerID = "post-to-discord-consumer"

	binanceClient *clients.BinanceClient

	coingeckoClient *coingecko.CoinGeckoClient
	coingeckoMtx    sync.Mutex

	discordClient           *clients.DiscordClient
	discordTwitterChannel   = clients.DiscordTwitterChannel
	discordWhaleChannel     = clients.DiscordWhaleChannel
	discordTradersChannel   = clients.DiscordTradersChannel
	discordNewsChannel      = clients.DiscordNewsChannel
	discordExchangesChannel = clients.DiscordExchangesChannel
	discordProjectsChannel  = clients.DiscordProjectsChannel
	discordTestingChannel   = clients.DiscordTestingChannel
	discordPriceBotChannel  = clients.DiscordPriceBotChannel

	defaultAlertsChannel  = clients.DiscordAlertsChannel
	defaultAlertsInterval = time.Duration(10 * time.Minute)

	alerterMap = map[string]*VolatilityAlerter{}
	alerterMtx sync.Mutex

	defaultATHAlertInterval = time.Duration(30 * time.Minute)
	athMap                  = map[string]*ATHAlerter{}
	athMtx                  sync.Mutex

	// Price Bot
	priceBot *PriceBot

	adminIDs = map[string]bool{
		"814142503393558558": true, // VinnyMac
		"805513165428883487": true, // ajperkins
	}

	insults = []string{
		"Come on, at least give me a ticker such as ETHUSDT.",
		"Mate, that is poggers. Give me a ticker like BTCUSDT",
		"Satoshi didn't build a blockchain for this. Ticker please.",
	}
)

func init() {
	register(PostToDiscordConsumer)

	ctx := context.Background()

	// Clients
	discordClient = clients.NewDiscordClient()
	binanceClient = clients.NewBinanceClient()
	coingeckoClient = coingecko.New(ctx)

	// Price Bot
	priceBot = NewPriceBot(ctx, discordPriceBotChannel, coingeckoClient, discordClient)
	go priceBot.Start(ctx)

	// Default ATH Alerting Setup
	athMtx.Lock()
	defer athMtx.Unlock()

	defaultATHCoins := GetDefaultATHAlertCoins()
	withJitter := len(defaultATHCoins) > 1
	slog.Info(context.TODO(), "Setting up default ATH alerts for %v", defaultATHCoins)
	for k := range defaultATHCoins {
		athAlerter := NewATHAlerter(k, defaultATHAlertInterval, discordClient, coingeckoClient, withJitter)
		athMap[k] = athAlerter
		go athAlerter.Run(context.TODO())
	}

	// Handlers
	discordClient.AddHandler(pingHandler)
	discordClient.AddHandler(priceHandler)
	discordClient.AddHandler(alerterHandler)
	discordClient.AddHandler(dealerterHandler)
	discordClient.AddHandler(athHandler)
	discordClient.AddHandler(whoIsThatHandler)
	discordClient.AddHandler(riskCalculator)

	// Check connectivity
	if err := binanceClient.Ping(); err != nil {
		slog.Error(nil, "Failed to connect to binance: %v", err)
	}

	slog.Info(nil, "Binance client connected")
}

func PostToDiscordConsumer(tweet *twitter.Tweet) {
	// Filter all tweets not made by pre-defined users.
	user, ok := GetMetadataMapping(tweet.User.ScreenName)
	if !ok {
		return
	}

	// Don't care for RT's
	if strings.HasPrefix(tweet.Text, "RT") {
		return
	}

	content := formatTwitterMessage(tweet)

	slog.Info(nil, content)

	if err := discordClient.PostToChannel(context.Background(), user.DiscordChannel, content); err != nil {
		slog.Error(nil, "Failed to post to discord: %v", err)
	}
}

func formatTwitterMessage(tweet *twitter.Tweet) string {
	m, ok := GetMetadataMapping(tweet.User.ScreenName)
	if !ok {
		return fmt.Sprintf("@%s [%v]: %s", tweet.User.ScreenName, tweet.CreatedAt, tweet.Text)
	}
	if len(m.Tags) == 0 {
		return fmt.Sprintf("@%s [%v]: %s", tweet.User.ScreenName, tweet.CreatedAt, tweet.Text)
	}

	roleIDs := []string{}
	for _, tag := range m.Tags {
		roleID, ok := GetRoleID(tag)
		if !ok {
			// Best effort
			slog.Info(context.TODO(), "Failed to find role id for tag: %s", tag)
			continue
		}
		roleIDs = append(roleIDs, roleID)
	}

	tagsStr := strings.Join(roleIDs, " ")
	return fmt.Sprintf("@%s [%v]: %s %s", tweet.User.ScreenName, tweet.CreatedAt, tweet.Text, tagsStr)
}

func pingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!ping") {
		return
	}

	slog.Info(nil, "Received PING from: %s", m.Author.Username)
	s.ChannelMessageSend(m.ChannelID, "PONG")
}

func priceHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!price") {
		return
	}

	symbols := strings.Split(m.Content, " ")

	if len(symbols) < 2 {
		s.ChannelMessageSend(m.ChannelID, randomInsultGenerator())
		return
	}

	if strings.ToLower(symbols[1]) == "all" {
		err := priceBot.PostLatestPrices(context.Background())
		switch {
		case err != nil:
			slog.Error(context.Background(), "Failed to post latest prices to discord: %v", err)
			return
		}
		slog.Trace(context.Background(), "Latest prices posted to discord")
		return
	}

	symbols = symbols[1:]

	for _, symbol := range symbols {
		slog.Info(nil, "Received !price cmd for %s from: %s", symbol, m.Author.Username)

		price, err := coingeckoClient.GetCurrentPriceFromSymbol(context.TODO(), symbol, "usd")
		if err != nil {
			slog.Error(nil, "Failed to get %s, err -> %v", symbol, err)
		}

		formattedPrice, err := util.FormatPriceAsString(price)
		if err != nil {
			slog.Info(nil, "Failedl to format price: %v", err.Error())
			s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("Sorry @%s, failed to retreive price for %s", m.Author.Username, symbol),
			)
			continue
		}

		var msg string
		switch symbol {
		case "BTCUSDT":
			msg = fmt.Sprintf("<:btc:816794855685750794> %s: %s", symbol, formattedPrice)
		case "ETHUSDT":
			msg = fmt.Sprintf("<:eth:817873038673707029> %s: %s", symbol, formattedPrice)
		default:
			msg = fmt.Sprintf("%s: %s", symbol, formattedPrice)
		}

		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func alerterHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!alerter") {
		return
	}

	symbols := strings.Split(m.Content, " ")
	if len(symbols) < 2 {
		s.ChannelMessageSend(m.ChannelID, randomInsultGenerator())
		return
	}

	alerterMtx.Lock()
	defer alerterMtx.Unlock()

	if symbols[1] == "GET" {
		if len(alerterMap) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No alerts are currently set, use `!alerter [<ticker>]` to set an alert.")
		}
		for symbol := range alerterMap {
			trigger := getTrigger(symbol)
			triggerPerc := fmt.Sprintf("%.1f", trigger*100)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Alert on: %s every %v +/- %s%%.", symbol, defaultAlertsInterval, triggerPerc))
		}
		return
	}

	symbols = symbols[1:]
	withJitter := len(symbols) > 1
	for _, symbol := range symbols {
		slog.Info(nil, "Received !alerter cmd for %s from: %s", symbol, m.Author.Username)
		a := NewVolatilityAlerter(symbol, binanceClient, discordClient, defaultAlertsChannel, defaultAlertsInterval, withJitter)

		if _, ok := alerterMap[symbol]; ok {
			continue
		}

		alerterMap[symbol] = a
		go a.Run(context.TODO())

		trigger := getTrigger(symbol)
		triggerPerc := fmt.Sprintf("%.1f", trigger*100)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Thanks @%s, alert registered for %s over %v  +/-%s%% .", m.Author.Username, symbol, defaultAlertsInterval, triggerPerc))
	}

}

func dealerterHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!dealerter") {
		return
	}

	symbols := strings.Split(m.Content, " ")

	if len(symbols) < 2 {
		s.ChannelMessageSend(m.ChannelID, randomInsultGenerator())
		return
	}
	alerterMtx.Lock()
	defer alerterMtx.Unlock()

	if symbols[1] == "ALL" {
		if len(alerterMap) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No alerts currently set.")
		}

		keys := []string{}
		for key := range alerterMap {
			keys = append(keys, key)
		}

		for _, key := range keys {
			if a, ok := alerterMap[key]; ok {
				a.Done()
				delete(alerterMap, key)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Thanks @%s, alert de-registered for %s every %v.", m.Author.Username, key, defaultAlertsInterval))
			}
		}
		return
	}

	symbols = symbols[1:]
	for _, symbol := range symbols {
		slog.Info(nil, "Received !dealerter cmd for %s from: %s", symbol, m.Author.Username)

		if a, ok := alerterMap[symbol]; ok {
			a.Done()

			delete(alerterMap, symbol)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Thanks @%s, alert de-registered for %s every %v", m.Author.Username, symbol, defaultAlertsInterval))
			continue
		}
	}
}

func athHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!ath") {
		return
	}

	slog.Info(context.TODO(), "!ath command received %s", m.Content)

	symbols := strings.Split(m.Content, " ")
	if len(symbols) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("@%s give me a single assest symbol eg `!ath BTC`", m.Author.Username))
		return
	}

	athMtx.Lock()
	defer athMtx.Unlock()

	if symbols[1] == "GET" {
		if len(athMap) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No alerts are currently set, use `!alerter [<ticker>]` to set an alert. NOTE: you don't need to incl. `USDT`")
		}
		for symbol := range athMap {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("ATH Alert on: %s", symbol))
		}
		return
	}

	symbols = symbols[1:]
	withJitter := len(symbols) > 1

	for _, symbol := range symbols {
		if _, ok := athMap[symbol]; ok {
			slog.Info(context.TODO(), "ATH Alerter already set %s; ignoring.", symbol)
			continue
		}

		athAlerter := NewATHAlerter(symbol, defaultATHAlertInterval, discordClient, coingeckoClient, withJitter)
		athMap[symbol] = athAlerter

		go athAlerter.Run(context.TODO())

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Thanks @%s, ATH alert recieved for %s", m.Author.Username, symbol))
	}
}

func whoIsThatHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!whoisthat") {
		return
	}

	tokens := strings.Split(m.Content, " ")

	if len(tokens) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("At least give me a twitter handle @%s...", m.Author.Username))
		return
	}

	twitterUser := tokens[1]
	twitterMetaData, ok := GetMetadataMapping(twitterUser)
	if !ok {
		slog.Info(nil, "No metadata stored for user: %v", twitterUser)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sorry @%s, no twitter metadata stored for: %s", m.Author.Username, twitterUser))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Username: %s\nName: %s\nBio: %s\nTwitter: %s\nYoutube: %s\nTwitch: %s\n", twitterUser, twitterMetaData.Name, twitterMetaData.Bio, twitterMetaData.Twitter, twitterMetaData.Youtube, twitterMetaData.Twitch))
}

func riskCalculator(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!risk") {
		return
	}

	tokens := strings.Split(m.Content, " ")
	if len(tokens) != 5 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, `!risk usage: <entry> <stop loss> <account size> <percentage eg 0.05>`", m.Author.Username))
		return
	}
	entry, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse entry: %v into a float, please check.", m.Author.Username, tokens[1]))
		return
	}
	stopLoss, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse stop loss: %v into a float, please check.", m.Author.Username, tokens[2]))
		return
	}
	accountSize, err := strconv.ParseFloat(tokens[3], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse accountSize: %v into a float, please check.", m.Author.Username, tokens[3]))
		return
	}
	percentage, err := strconv.ParseFloat(tokens[4], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse percentage: %v into a float, please check.", m.Author.Username, tokens[4]))
		return
	}

	contracts := calculateRisk(entry, stopLoss, accountSize, percentage)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, you need to buy **%.2f** contracts for %v%% risk.", m.Author.Username, contracts, percentage*100))
	return
}

func randomInsultGenerator() string {
	nIndexes := len(insults)
	return insults[rand.Intn(nIndexes)]
}

// calculateRisk returns the number of contracts to buy.
func calculateRisk(entry, stopLoss, accountSize, percentage float64) float64 {
	maxRiskToLose := percentage * accountSize
	lossPerContract := entry - stopLoss
	return maxRiskToLose / lossPerContract
}
