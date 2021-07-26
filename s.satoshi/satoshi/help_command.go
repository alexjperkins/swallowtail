package satoshi

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	helpCommandID     = "help-command"
	helpCommandPrefix = "!help"
)

func init() {
	registerSatoshiCommand(helpCommandID, helpCommand)
}

func helpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, helpCommandPrefix) {
		return
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: Hi <@%s>, this command isn't yet implemented :disappointed:", m.Author.ID),
	)
}
