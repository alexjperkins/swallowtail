package clients

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"swallowtail/libraries/util"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

var (
	discordToken       string
	discordBotUsername = "satoshi"
	discordGuildID     = "814144801977008228"

	discordMtx sync.Mutex

	// Channels
	DiscordTwitterChannel   = "twitter"
	DiscordAlertsChannel    = "bot-alerts"
	DiscordTestingChannel   = "testing"
	DiscordWhaleChannel     = "whale-alerts"
	DiscordTradersChannel   = "traders-feed"
	DiscordNewsChannel      = "crypto-news"
	DiscordExchangesChannel = "exchanges-alerts"
	DiscordProjectsChannel  = "project-alerts"
	DiscordPriceBotChannel  = "price-bot"

	channelIDMapping = map[string]string{
		DiscordTwitterChannel:   "816794087868465163",
		DiscordAlertsChannel:    "816794120851816479",
		DiscordTestingChannel:   "817513133274824715",
		DiscordWhaleChannel:     "817789196319195166",
		DiscordTradersChannel:   "817789261415448606",
		DiscordNewsChannel:      "817789219656826970",
		DiscordExchangesChannel: "818909423530541116",
		DiscordProjectsChannel:  "826528849374216192",
		DiscordPriceBotChannel:  "831234720943702066",
	}
)

func init() {
	discordMtx.Lock()
	defer discordMtx.Unlock()
	discordToken = util.SetEnv("DISCORD_API_TOKEN")
}

type Attachment struct {
	link string
}

// New creates a new discord client
func NewDiscordClient() *DiscordClient {
	s, err := discordgo.New(fmt.Sprintf("Bot %s", discordToken))
	if err != nil {
		panic(err)
	}

	done := make(chan struct{}, 1)
	err = s.Open()
	if err != nil {
		slog.Error(nil, "Could not open discord ws")
	}

	slog.Info(nil, "discord ws opened")
	go func() {
		defer slog.Info(nil, "Closing down discord ws")
		defer s.Close()

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		select {
		case <-sc:
			return
		case <-done:
			return
		}
	}()

	slog.Info(nil, "Created discord bot: %s", discordBotUsername)

	return &DiscordClient{
		session: s,
		done:    done,
	}
}

type DiscordClient struct {
	session *discordgo.Session
	done    chan struct{}
}

func (d *DiscordClient) PostToChannel(ctx context.Context, channel, content string) error {
	channelID, ok := channelIDMapping[channel]
	if !ok {
		return fmt.Errorf("Failed to find channel id for channel %s", channel)
	}

	msg, err := d.session.ChannelMessageSend(channelID, content)
	if err != nil {
		return err
	}
	slog.Info(ctx, "Message Posted to discord: %v", msg)
	return nil
}

// AddHandler adds handler to discord that listens to events
func (d *DiscordClient) AddHandler(handler interface{}) func() {
	return d.session.AddHandler(handler)
}

func (d *DiscordClient) BidirectionalSubscription(ctx context.Context) error {
	if err := d.session.Open(); err != nil {
		return err
	}

	defer d.session.Close()
	sc := make(chan os.Signal, 1)
	select {
	case <-sc:
		return nil
	case <-ctx.Done():
		return nil
	}
}

func (d *DiscordClient) Stop() {
	select {
	case d.done <- struct{}{}:
		return
	}
}
