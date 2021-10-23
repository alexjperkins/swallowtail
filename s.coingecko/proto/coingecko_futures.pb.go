package coingeckoproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Get Asset Latest Price By Symbol --- //
type GetAssetLatestPriceBySymbolFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *GetAssetLatestPriceBySymbolResponse
	ctx     context.Context
}

func (a *GetAssetLatestPriceBySymbolFuture) Response() (*GetAssetLatestPriceBySymbolResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_latest_asset_price_by_symbol", err)
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

func (r *GetAssetLatestPriceBySymbolRequest) Send(ctx context.Context) *GetAssetLatestPriceBySymbolFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetAssetLatestPriceBySymbolRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetAssetLatestPriceBySymbolFuture {
	errc := make(chan error, 1)
	resultc := make(chan *GetAssetLatestPriceBySymbolResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-coingecko:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s-coingecko_connection_failed", nil)
		return &GetAssetLatestPriceBySymbolFuture{
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
	c := NewCoingeckoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetAssetLatestPriceBySymbol(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_get_latest_asset_price_by_symbol_grpc_call", nil)
			return
		}
		resultc <- rsp
	}()

	return &GetAssetLatestPriceBySymbolFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Get ATH By Symbol --- //
type GetATHBySymbolFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *GetATHBySymbolResponse
	ctx     context.Context
}

func (a *GetATHBySymbolFuture) Response() (*GetATHBySymbolResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_ath_by_symbol", err)
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

func (r *GetATHBySymbolRequest) Send(ctx context.Context) *GetATHBySymbolFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetATHBySymbolRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetATHBySymbolFuture {
	errc := make(chan error, 1)
	resultc := make(chan *GetATHBySymbolResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-coingecko:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s-coingecko_connection_failed", nil)
		return &GetATHBySymbolFuture{
			ctx:  ctx,
			errc: errc,
			closer: func() error {
				if conn != nil {
					conn.Close()
				}
				return nil
			},
			resultc: resultc,
		}
	}
	c := NewCoingeckoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetATHBySymbol(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_get_ath_by_symbol_grpc_call", nil)
			return
		}
		resultc <- rsp
	}()

	return &GetATHBySymbolFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
