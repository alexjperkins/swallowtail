package tradeengineproto

import (
	context "context"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"

	"swallowtail/libraries/gerrors"
)

// --- Create Trade Strategy --- //

type CreateTradeStrategyFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *CreateTradeStrategyResponse
	ctx     context.Context
}

func (a *CreateTradeStrategyFuture) Response() (*CreateTradeStrategyResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "create_trade_strategy", err)
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

func (r *CreateTradeStrategyRequest) Send(ctx context.Context) *CreateTradeStrategyFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *CreateTradeStrategyRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *CreateTradeStrategyFuture {
	errc := make(chan error, 1)
	resultc := make(chan *CreateTradeStrategyResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-tradeengine:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &CreateTradeStrategyFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewTradeengineClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.CreateTradeStrategy(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_create_trade_strategy", nil)
			return
		}
		resultc <- rsp
	}()

	return &CreateTradeStrategyFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Read Trade --- //

type ReadTradeByTradeIDFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ReadTradeByTradeIDResponse
	ctx     context.Context
}

func (a *ReadTradeByTradeIDFuture) Response() (*ReadTradeByTradeIDResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_trade_by_trade_id", err)
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

func (r *ReadTradeByTradeIDRequest) Send(ctx context.Context) *ReadTradeByTradeIDFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ReadTradeByTradeIDRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ReadTradeByTradeIDFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ReadTradeByTradeIDResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-tradeengine:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ReadTradeByTradeIDFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewTradeengineClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ReadTradeByTradeID(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_read_trade_by_trade_id", nil)
			return
		}
		resultc <- rsp
	}()

	return &ReadTradeByTradeIDFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Execute Trade Strategy For Participant --- //

type ExecuteTradeStrategyForParticipantFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ExecuteTradeStrategyForParticipantResponse
	ctx     context.Context
}

func (a *ExecuteTradeStrategyForParticipantFuture) Response() (*ExecuteTradeStrategyForParticipantResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "execute_trade_strategy_for_participant", err)
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

func (r *ExecuteTradeStrategyForParticipantRequest) Send(ctx context.Context) *ExecuteTradeStrategyForParticipantFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ExecuteTradeStrategyForParticipantRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ExecuteTradeStrategyForParticipantFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ExecuteTradeStrategyForParticipantResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-tradeengine:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ExecuteTradeStrategyForParticipantFuture{
			ctx:  ctx,
			errc: errc,
			closer: func() error {
				if conn != nil {
					return conn.Close()
				}
				return nil
			},
			resultc: resultc,
		}
	}
	c := NewTradeengineClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ExecuteTradeStrategyForParticipant(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_execute_trade_strategy_for_trade_participant", nil)
			return
		}
		resultc <- rsp
	}()

	return &ExecuteTradeStrategyForParticipantFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
