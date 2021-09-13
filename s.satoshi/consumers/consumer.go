package consumers

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ConsumerMessage the struct definition of a message that satoshi will consume.
// One caveat is that we are vendor locked with Discord with the name conventions &
// use of attachments here - whilst practically this won't change, it would be nice to
// be agnostic.
type ConsumerMessage struct {
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
	// Poller is a function that is called that monitors a given message over some period.
	// This is will be called synchronously by the Satoshi client stream if a message contains it.
	Poller func(ctx context.Context, messageID string) error
}

// Consumer ...
type Consumer interface {
	// Receiver accepts a channel with the intention of asynchronously publishing messages too it.
	Receiver(ctx context.Context, c chan *ConsumerMessage, d chan struct{}, withJitter bool)
	// IsActive simply returns a boolean that indicates if the consumer is active or not.
	// The premise here, is that we have the concept of "shadow consumers"; which consume messages,
	// but shouldn't do anything with them other than log.
	IsActive() bool
}
