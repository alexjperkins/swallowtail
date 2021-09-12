package tradeengineproto

import (
	context "context"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"

	"swallowtail/libraries/gerrors"
)

// --- Create Trade --- //

type CreateTradeFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *CreateTradeResponse
	ctx     context.Context
}

func (a *CreateTradeFuture) Response() (*CreateTradeResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "create_trade", err)
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

func (r *CreateTradeRequest) Send(ctx context.Context) *CreateTradeFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *CreateTradeRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *CreateTradeFuture {
	errc := make(chan error, 1)
	resultc := make(chan *CreateTradeResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-tradeengine:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &CreateTradeFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewTradeengineClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.CreateTrade(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_create_trade", nil)
			return
		}
		resultc <- rsp
	}()

	return &CreateTradeFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Create Trade --- //

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
