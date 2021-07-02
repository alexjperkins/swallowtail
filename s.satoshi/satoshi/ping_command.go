package satoshi

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	pingCommandID = "pingCommand"

	pingCommandPrefix = "!ping"
)

func init() {
	registerSatoshiCommand(pingCommandID, pingCommand)
}

func pingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, pingCommandPrefix) {
		return
	}

	slog.Info(context.TODO(), "Received %s command from: %s", pingCommandPrefix, m.Author.Username)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: Hello <@%s>, what can I do for you?", m.Author.ID))
}
