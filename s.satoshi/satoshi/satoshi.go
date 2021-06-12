package satoshi

import (
	"context"
	"fmt"
	"swallowtail/libraries/util"
	discord "swallowtail/s.discord/clients"

	"github.com/monzo/slog"
)

var (
	satoshiToken = util.SetEnv("DISCORD_API_TOKEN")
	SatoshiBotID = "satoshi-bot"
)

// Satoshi Interface
type Satoshi interface {
	// Run simply starts Satohsi & all registered consumers & command handlers.
	Run(ctx context.Context)
	// Stop stops satoshi gracefully.
	Stop()
}

func New(withJitter bool) Satoshi {
	dc := discord.New(SatoshiBotID, satoshiToken, true)
	for id, command := range commandRegistry {
		slog.Info(context.TODO(), "Registering command %s to %s", id, SatoshiBotID)
		dc.AddHandler(command)
	}
	return &satoshi{
		dc:             dc,
		withJitter:     withJitter,
		consumers:      consumerRegistry,
		consumerStream: make(chan *SatoshiConsumerMessage, 32),
		done:           make(chan struct{}, 1),
	}
}

type satoshi struct {
	dc             discord.DiscordClient
	withJitter     bool
	consumers      map[string]SatoshiConsumer
	consumerStream chan *SatoshiConsumerMessage
	done           chan struct{}
}

func (s *satoshi) Run(ctx context.Context) {
	s.consume(ctx)
	go s.streamEventHandler(ctx)
}

func (s *satoshi) Stop() {
	slog.Info(context.TODO(), "Satoshi stop signal received.")
	defer close(s.consumerStream)
	select {
	case s.done <- struct{}{}:
	default:
		slog.Warn(context.TODO(), "Cannot stop satoshi; blocked done channel")
	}
}

func (s *satoshi) consume(ctx context.Context) {
	for id, c := range s.consumers {
		id, c := id, c
		go func() {
			slog.Info(ctx, "Starting registered satoshi consumer %s", id)
			c.Receiver(ctx, s.consumerStream, s.done, s.withJitter)
		}()
	}
}

func (s *satoshi) streamEventHandler(ctx context.Context) {
	for {
		select {
		case e := <-s.consumerStream:
			slog.Info(ctx, "Received satoshi consumer event; %+v, isActive: %b", e, e.IsActive)
			if !e.IsActive {
				continue
			}
			if e.IsPrivate {
				if len(e.ParticipentIDs) == 0 {
					slog.Warn(
						ctx, "Dropping event; cannot send private message with no participents.",
						map[string]string{
							"event": fmt.Sprintf("%+v", e),
						},
					)
					continue
				}
				// Currently we only have the functionality to send to one participent.
				participent := e.ParticipentIDs[0]
				s.dc.SendPrivateMessage(ctx, e.Message, participent)
				continue
			}
			// Send the event message to the appropriate discord channel.
			s.dc.Send(ctx, e.Message, e.DiscordChannelID)
		case <-ctx.Done():
			return
		case <-s.done:
			return
		}
	}
}
