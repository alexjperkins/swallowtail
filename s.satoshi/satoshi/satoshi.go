package satoshi

import (
	"context"
	"fmt"

	"github.com/monzo/slog"

	"swallowtail/libraries/util"
	discord "swallowtail/s.discord/client"
	"swallowtail/s.satoshi/commands"
	"swallowtail/s.satoshi/consumers"
)

const (
	Version = "v0.1.9"
)

var (
	satoshiToken = util.SetEnv("DISCORD_API_TOKEN")
	SatoshiBotID = "satoshi-bot"
)

// Initializes satoshi background processes.
func Init(ctx context.Context) error {
	dc := discord.New(SatoshiBotID, satoshiToken, true)

	for id, command := range commands.List() {
		slog.Info(context.TODO(), "Registering command %d) %s to %s", id, command.ID, SatoshiBotID)
		dc.AddHandler(command.Exec)
	}

	s := &satoshi{
		dc:             dc,
		withJitter:     true,
		consumers:      consumers.Registry(),
		consumerStream: make(chan *consumers.ConsumerMessage, 32),
		done:           make(chan struct{}, 1),
	}

	s.run(ctx)
	return nil
}

type satoshi struct {
	dc             discord.DiscordClient
	withJitter     bool
	consumers      map[string]consumers.Consumer
	consumerStream chan *consumers.ConsumerMessage
	done           chan struct{}
}

func (s *satoshi) run(ctx context.Context) {
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
