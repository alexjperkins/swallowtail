package satoshiproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Poll Trade Participants--- //

type PollTradeParticipantsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *PollTradeParticipantsResponse
	ctx     context.Context
}

func (a *PollTradeParticipantsFuture) Response() (*PollTradeParticipantsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "poll_trade_participants", err)
		}
	}()

	select {
	case r := <-a.resultc:
		return r, nil
	case <-a.ctx.Done():
		return nil, a.ctx.Err()
	case err := <-a.errc:
		return nil, err
	}
}

func (r *PollTradeParticipantsRequest) Send(ctx context.Context) *PollTradeParticipantsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *PollTradeParticipantsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *PollTradeParticipantsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *PollTradeParticipantsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-satoshi:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &PollTradeParticipantsFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewSatoshiClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.PollTradeParticipants(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_poll_trade_participants", nil)
			return
		}
		resultc <- rsp
	}()

	return &PollTradeParticipantsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
