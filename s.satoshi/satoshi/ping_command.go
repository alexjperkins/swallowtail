package satoshi

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	pingCommandID = "pingCommand"
)

func init() {
	registerSatoshiCommand(pingCommandID, pingCommand)
}

func pingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!ping") {
		return
	}

	slog.Info(nil, "Received PING from: %s", m.Author.Username)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: Hello <@%s>, what can I do for you?", m.Author.ID))
}
