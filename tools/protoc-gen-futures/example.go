package coingeckoproto

import (
	context "context"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- GetAssetLatestPriceByID --- //
type GetAssetLatestPriceByIDFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *GetAssetLatestPriceByIDResponse
	ctx     context.Context
}

func (a *GetAssetLatestPriceByIDFuture) Response() (*GetAssetLatestPriceByIDResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_latest_asset_price_by_id", err)
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

func (r *GetAssetLatestPriceByIDRequest) Send(ctx context.Context) *GetAssetLatestPriceByIDFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetAssetLatestPriceByIDRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetAssetLatestPriceByIDFuture {
	errc := make(chan error, 1)
	resultc := make(chan *GetAssetLatestPriceByIDResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-coingecko:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &GetAssetLatestPriceByIDFuture{
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewCoingeckoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		defer func() {
			cancel()
		}()

		rsp, err := c.GetAssetLatestPriceByID(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &GetAssetLatestPriceByIDFuture{
		closer: conn.Close,
		errc:   errc,
		ctx:    ctx,
	}
}

// --- GetAssetLatestPriceBySymbol --- //
type GetAssetLatestPriceBySymbolFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *GetAssetLatestPriceBySymbolResponse
	ctx     context.Context
}

func (a *GetAssetLatestPriceBySymbolFuture) Response() (*GetAssetLatestPriceBySymbolResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_latest_asset_price_by_id", err)
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
		errc <- err
		return &GetAssetLatestPriceBySymbolFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewCoingeckoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetAssetLatestPriceBySymbol(ctx, r)
		if err != nil {
			errc <- err
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
