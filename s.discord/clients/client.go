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

// New creates a new discord client
func New() *DiscordClient {
	s, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		panic(err)
	}

	slog.Info(nil, "Created discord bot: %s", discordBotUsername)

	return &DiscordClient{
		session: s,
	}
}

type Handler func(context.Context, *discordgo.Session, *discordgo.MessageCreate) error

type Attachment struct {
	link string
}

type DiscordClient struct {
	session *discordgo.Session
}

func (d *DiscordClient) Send(ctx context.Context, channel, content string) error {
	msg, err := d.session.ChannelMessageSend(channel, content)
	if err != nil {
		return err
	}
	slog.Info(ctx, "Message Posted to discord: %v", msg)
	return nil
}

func (d *DiscordClient) SendPrivateMessage(ctx context.Context, discordID string, content string) error {
	channelID, ok := privateChannelMapping[discordID]
	if ok {
		return d.Send(ctx, channelID, content)
	}

	ch, err := d.session.UserChannelCreate(discordID)
	if err != nil {
		return terrors.Augment(err, "Failed to create private channel", map[string]string{
			"discord_user_id": discordID,
		})
	}
	return d.Send(ctx, ch.ID, content)
}

func (d *DiscordClient) AddHandlerToAllChannel(ctx context.Context, handler Handler) error {
	return nil
}

func (d *DiscordClient) AddHandlerToChannel(ctx context.Context, channel string, handler Handler) func() {
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
