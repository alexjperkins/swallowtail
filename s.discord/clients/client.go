package clients

import (
	"context"
	"os"
	"swallowtail/libraries/util"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

var (
	DiscordClientID    = "discord-client-id"
	discordToken       string
	discordBotUsername = "satoshi"

	mtx sync.Mutex

	discordTwitterChannel = "twitter"
	channelIDMapping      = map[string]string{
		discordTwitterChannel: "816794087868465163",
	}

	privateChannelMapping = map[string]string{}
	pMu                   sync.Mutex
)

func init() {
	mtx.Lock()
	defer mtx.Unlock()
	discordToken = util.SetEnv("DISCORD_API_TOKEN")
}

type DiscordClient interface {
	Send(ctx context.Context, message, channelID string) error
	SendPrivateMessage(ctx context.Context, message, userID string) error
	Subscribe(ctx context.Context, subscriberID string) error
}

// New creates a new discord client
func New() DiscordClient {
	s, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		panic(terrors.Augment(err, "Failed to create discord client", nil))
	}
	slog.Info(nil, "Created discord bot: %s", discordBotUsername)
	return &discordClient{
		session: s,
	}
}

type Handler func(context.Context, *discordgo.Session, *discordgo.MessageCreate) error

type discordClient struct {
	session *discordgo.Session
}

func (d *discordClient) Send(ctx context.Context, message, channelID string) error {
	msg, err := d.session.ChannelMessageSend(channelID, message)
	if err != nil {
		return err
	}
	slog.Info(ctx, "Message Posted to discord: %v", msg)
	return nil
}

func (d *discordClient) SendPrivateMessage(ctx context.Context, message, userID string) error {
	channelID, ok := privateChannelMapping[userID]
	if ok {
		return d.Send(ctx, channelID, message)
	}
	ch, err := d.session.UserChannelCreate(message)
	if err != nil {
		return terrors.Augment(err, "Failed to create private channel", map[string]string{
			"discord_user_id": userID,
		})
	}
	return d.Send(ctx, ch.ID, message)
}

func (d *discordClient) Subscribe(ctx context.Context, subscriberID string) error {
	return nil
}

func (d *discordClient) AddHandlerToChannel(ctx context.Context, channel string, handler Handler) func() {
	return d.session.AddHandler(handler)
}

func (d *discordClient) BidirectionalSubscription(ctx context.Context) error {
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
