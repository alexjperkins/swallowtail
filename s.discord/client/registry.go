package client

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// GuildMemberAddHandler ...
type GuildMemberAddHandler func(s *discordgo.Session, u *discordgo.GuildMemberAdd)

var (
	guildMemberAddRegistry = map[string]GuildMemberAddHandler{}
	guildMemberAddMu       sync.RWMutex
)

func registerGuildMemberAddHandler(id string, handler GuildMemberAddHandler) {
	guildMemberAddMu.Lock()
	defer guildMemberAddMu.Unlock()

	if _, ok := guildMemberAddRegistry[id]; ok {
		log.Fatalf("Failed to register guild member add handler: %s: duplicate id", id)
	}

	guildMemberAddRegistry[id] = handler
}
