package satoshi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"swallowtail/libraries/util"
	binanceclient "swallowtail/s.binance/client"
	discordclient "swallowtail/s.discord/client"
	discordproto "swallowtail/s.discord/proto"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

const (
	discordConsumerID = "discord-consumer"
	discordChannelUrl = "https://discord.com/api/v9/channels/%s/messages?limit=1"
)

var (
	discordConsumerToken string
	discordConsumerName  = "satoshi_consumer"

	binanceAssetPairs = map[string]bool{}
)

func init() {
	discordConsumerToken = util.SetEnv("SATOSHI_DISCORD_CONSUMER_1_API_TOKEN")
	registerSatoshiConsumer(discordConsumerID, DiscordConsumer{
		Active: true,
	})

}

type DiscordConsumer struct {
	Active bool
}

func (dc DiscordConsumer) Receiver(ctx context.Context, c chan *SatoshiConsumerMessage, d chan struct{}, _ bool) {
	discordClient := discordclient.New(discordConsumerName, discordConsumerToken, false)

	// Build a list of asset pairs that have futures trading enabled.
	assets, err := binanceclient.ListAllAssetPairs(context.Background())
	if err != nil {
		panic(err)
	}
	for _, asset := range assets.Symbols {
		if asset.WithMarginTrading {
			binanceAssetPairs[strings.ToLower(asset.BaseAsset)] = true
		}
	}
	slog.Info(ctx, "Fetched all binance asset pairs for the discord consumer; total: %v", len(binanceAssetPairs))

	// Add handlers
	discordClient.AddHandler(handleModMessages(ctx, c, dc.Active))
	discordClient.AddHandler(handleSwingMessages(ctx, c, dc.Active))

	defer slog.Warn(ctx, "Discord consumer stop signal received.")
	defer discordClient.Close()

	select {
	case <-d:
	case <-ctx.Done():
	}
}

func (dc DiscordConsumer) IsActive() bool {
	return dc.Active
}

func formatContent(ctx context.Context, username, timestamp, content string) string {
	ts, err := time.Parse(time.RFC3339, timestamp)
	switch {
	case err != nil:
		slog.Warn(ctx, "Failed to parse timestamp; setting as original: %s, err: %v", timestamp, err)
	default:
		timestamp = ts.Truncate(time.Minute).String()
	}
	return fmt.Sprintf("%s[%v]:\n```%s```", username, timestamp, content)
}

func handleModMessages(
	ctx context.Context, c chan *SatoshiConsumerMessage, isActive bool,
) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, mc *discordgo.MessageCreate) {
		m := mc.Message
		if m.ChannelID != discordproto.DiscordMoonModMessagesChannel {
			return
		}
		parsedContent, err := getLatestChannelMessages(ctx, s, m.ChannelID)
		if err != nil {
			slog.Error(ctx, "Failed to get latest mod message: %v", err)
			return
		}
		msgs := []*SatoshiConsumerMessage{}
		for i, pc := range parsedContent {
			// First lets check if the content is part of any 1-10k challenge.
			if !contains1To10kChallenge(strings.ToLower(pc.Content)) {
				slog.Debug(ctx, "1-10k challenge message received: %s", pc.Content)
				msg := &SatoshiConsumerMessage{
					ConsumerID:       discordConsumerID,
					DiscordChannelID: discordproto.DiscordSatoshiGeneralChannel,
					Message:          warning("1-10k challenge update from Astekz", pc.Content),
					Created:          time.Now(),
					IsActive:         isActive,
				}
				msgs = append(msgs, msg)
			}

			// We want to reduce the noise here; we only care if the content contains a Ticker.
			if !containsTicker(strings.ToLower(pc.Content)) {
				slog.Debug(ctx, "Received mod message without ticker: %s", pc.Content)
				return
			}
			msg := &SatoshiConsumerMessage{
				ConsumerID:       discordConsumerID,
				DiscordChannelID: discordproto.DiscordSatoshiModMessagesChannel,
				Message:          formatContent(ctx, pc.Author.Username, pc.Timestamp, pc.Content),
				Created:          time.Now(),
				IsActive:         isActive,
				Metadata: map[string]string{
					"message": fmt.Sprintf("%v", i),
					"total":   fmt.Sprintf("%v", len(parsedContent)),
				},
			}
			msgs = append(msgs, msg)
		}

		// Lets publish our messages.
		for _, msg := range msgs {
			select {
			case c <- msg:
			default:
				slog.Warn(ctx, "Failed to publish satoshi mods msg; blocked channel")
			}

		}
	}
}

func handleSwingMessages(
	ctx context.Context, c chan *SatoshiConsumerMessage, isActive bool,
) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, mc *discordgo.MessageCreate) {
		m := mc.Message
		if m.ChannelID != discordproto.DiscordMoonSwingGroupChannel {
			return
		}
		parsedContent, err := getLatestChannelMessages(ctx, s, m.ChannelID)
		if err != nil {
			slog.Error(ctx, "Failed to get latest mod message: %v", err)
			return
		}

		for i, pc := range parsedContent {
			msg := &SatoshiConsumerMessage{
				ConsumerID:       discordConsumerID,
				DiscordChannelID: discordproto.DiscordSatoshiSwingsChannel,
				Message:          formatContent(ctx, pc.Author.Username, pc.Timestamp, pc.Content),
				Created:          time.Now(),
				IsActive:         isActive,
				Metadata: map[string]string{
					"message": fmt.Sprintf("%v", i),
					"total":   fmt.Sprintf("%v", len(parsedContent)),
				},
			}
			select {
			case c <- msg:
			default:
				slog.Warn(ctx, "Failed to publish satoshi swings msg; blocked channel")
			}
		}
	}
}

type channelMessage struct {
	ID        string                `json:"id"`
	ChannelID string                `json:"channel_id"`
	Author    *channelMessageAuthor `json:"author"`
	Content   string                `json:"content"`
	Timestamp string                `json:"timestamp"`
}

type channelMessageAuthor struct {
	Username string `json:"username"`
}

func getLatestChannelMessages(ctx context.Context, s *discordgo.Session, channelID string) ([]*channelMessage, error) {
	// TODO: A hack for testing purposes.
	if s == nil {
		return nil, nil
	}
	url := fmt.Sprintf(discordChannelUrl, channelID)
	slogParams := map[string]string{
		"channel_id": channelID,
		"url":        url,
	}
	// Create request
	req, err := http.NewRequestWithContext(
		ctx, "GET", url, nil,
	)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to create request", slogParams)
	}
	req.Header.Set("authorization", discordConsumerToken)

	// Execute request.
	rsp, err := s.Client.Do(req)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to make request", slogParams)
	}
	defer rsp.Body.Close()

	// Parse Body.
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to read response.", slogParams)
	}
	var msgList []*channelMessage
	json.Unmarshal([]byte(body), &msgList)

	return msgList, nil
}

func contains1To10kChallenge(content string) bool {
	var (
		contains1k     bool
		contains10k    bool
		containsAstekz bool
	)

	tokens := strings.Fields(content)
	for _, token := range tokens {
		token := strings.ToLower(token)
		if strings.Contains(token, "1k") {
			contains1k = true
		}

		if strings.Contains(token, "10k") {
			contains10k = true
		}

		if strings.Contains(token, "astekz") {
			containsAstekz = true
		}
	}

	return contains1k && contains10k && containsAstekz
}

// containsTicker checks if the contain contains a ticker that is traded on Binance
// it assumes that the content passed with be normalized to lowercase.
func containsTicker(content string) bool {
	tokens := strings.Fields(strings.ToLower(content))
	for _, token := range tokens {
		switch {
		case
			token == "usd",
			token == "usdt",
			token == "usdc":
			// If we match against some stablecoin inadvertly; then we can skip.
			continue
		case
			strings.Contains(token, "usd"),
			strings.Contains(token, "usdc"),
			strings.Contains(token, "usdt"):
			// But if a token contains a stable coins, then lets assume it's of the form BTCUSDT.
			// We might pick up typos and similar here, but that's fine for now.

			fmt.Println(token)

			return true
		case strings.Contains(token, "/"):
			// Some mods format their trades as `BTC/USDT`.
			childContent := strings.ReplaceAll(token, "/", " ")
			if containsTicker(childContent) {
				return true
			}
		}

		if _, ok := binanceAssetPairs[token]; ok {
			return true
		}
	}
	return false
}

// Formats a message for a standardized warning.
func warning(greeting, content string) string {
	return fmt.Sprintf(":rotating_light: %s :rotating_light:\n```%s```", greeting, content)
}
