package satoshi

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// SatoshiConsumerMessage the struct definition of a message that satoshi will consume.
// One caveat is that we are vendor locked with Discord with the name conventions &
// use of attachments here - whilst practically this won't change, it would be nice to
// be agnostic.
type SatoshiConsumerMessage struct {
	// The message attachments of which to post along with the content.
	Attachments []*discordgo.MessageAttachment
	// The ID of the consumer (sender).
	ConsumerID string
	// The ID of the discord channel of which to post to.
	DiscordChannelID string
	// The content of the message itself.
	Message string
	// An idempotency key used to prevent duplicates; for now this is for future-proofing.
	IdempotencyKey string
	// A boolean flag to indicate if it should be sent to a private channel.
	IsPrivate bool
	// To be used in conjunction with `IsPrivate`, the participients of whom to create a channel with.
	// A single participent indicates it's a direct message.
	ParticipentIDs []string
	// Basic metadata.
	Metadata map[string]string
	// The timestamp of when the message was created.
	Created time.Time
	// If this is true; then we don't send, rather we log.
	IsActive bool
}

// SatoshiCommand type definition of what a command is; alias for a discord handler.
type SatoshiCommand func(s *discordgo.Session, m *discordgo.MessageCreate)

// SatoshiConsumer is the contract that is required to be satisfied to add to satoshi
type SatoshiConsumer interface {
	// Receiver accepts a channel with the intention of asynchronously publishing messages too it.
	Receiver(ctx context.Context, c chan *SatoshiConsumerMessage, d chan struct{}, withJitter bool)
	// IsActive simply returns a boolean that indicates if the consumer is active or not.
	// The premise here, is that we have the concept of "shadow consumers"; which consume messages,
	// but shouldn't do anything with them other than log.
	IsActive() bool
}

var (
	commandRegistry  = map[string]SatoshiCommand{}
	consumerRegistry = map[string]SatoshiConsumer{}

	consumerMtx sync.Mutex
	commandMtx  sync.Mutex
)

// registerSatoshiCommand registers a new command with the given ID to Satoshi
func registerSatoshiCommand(id string, command SatoshiCommand) {
	commandMtx.Lock()
	defer commandMtx.Unlock()
	if _, ok := commandRegistry[id]; ok {
		panic(fmt.Sprintf("Can't register the same command twice: %s", id))
	}
	commandRegistry[id] = command
}

// registerSatoshiConsumer registers a new consumer with the given ID to Satoshi
func registerSatoshiConsumer(id string, consumer SatoshiConsumer) {
	consumerMtx.Lock()
	defer consumerMtx.Unlock()
	if _, ok := consumerRegistry[id]; ok {
		panic(fmt.Sprintf("Can't register the same consumer twice: %s", id))
	}
	consumerRegistry[id] = consumer
}
