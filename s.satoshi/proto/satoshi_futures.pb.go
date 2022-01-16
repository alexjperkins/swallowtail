package satoshiproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Poll Trade Strategy Participants--- //

type PollTradeStrategyParticipantsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *PollTradeStrategyParticipantsResponse
	ctx     context.Context
}

func (a *PollTradeStrategyParticipantsFuture) Response() (*PollTradeStrategyParticipantsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "poll_trade_strategy_participants", err)
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

func (r *PollTradeStrategyParticipantsRequest) Send(ctx context.Context) *PollTradeStrategyParticipantsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *PollTradeStrategyParticipantsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *PollTradeStrategyParticipantsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *PollTradeStrategyParticipantsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-satoshi:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &PollTradeStrategyParticipantsFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewSatoshiClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.PollTradeStrategyParticipants(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_poll_trade_strategy_participants", nil)
			return
		}
		resultc <- rsp
	}()

	return &PollTradeStrategyParticipantsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
