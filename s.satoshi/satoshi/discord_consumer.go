package satoshi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"swallowtail/libraries/util"
	discord "swallowtail/s.discord/clients"
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
	discordClient := discord.New(discordConsumerName, discordConsumerToken, false)

	// Add handlers
	discordClient.AddHandler(handleModMessages(ctx, c, dc.Active))
	discordClient.AddHandler(handleSwingMessages(ctx, c, dc.Active))

	defer slog.Warn(ctx, "Discord consumer stop signal received.")
	defer discordClient.Close()

	select {
	case <-d:
		return
	case <-ctx.Done():
		return
	}
}

func (dc DiscordConsumer) IsActive() bool {
	return dc.Active
}

func formatContent(username, timestamp, content string) string {
	return fmt.Sprintf("%s: [%v] %s", username, timestamp, content)
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
		for i, pc := range parsedContent {
			msg := &SatoshiConsumerMessage{
				ConsumerID:       discordConsumerID,
				DiscordChannelID: discordproto.DiscordSatoshiModMessagesChannel,
				Message:          formatContent(pc.Author.Username, pc.Timestamp, pc.Content),
				Created:          time.Now(),
				IsActive:         isActive,
			}
			select {
			case c <- msg:
				slog.Info(ctx, "Published mod msg %v/%v, %+v", i, len(parsedContent), msg)
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
				Message:          formatContent(pc.Author.Username, pc.Timestamp, pc.Content),
				Created:          time.Now(),
				IsActive:         isActive,
			}
			select {
			case c <- msg:
				slog.Info(ctx, "Published swing msg %v/%v, %+v", i, len(parsedContent), msg)
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
		return nil, terrors.Augment(err, "Oh no, failed to get response", slogParams)
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
