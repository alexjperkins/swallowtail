package binanceproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Execute Futures Perpetuals Trade --- //

type ExecuteFuturesPerpetualsTradeFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ExecuteFuturesPerpetualsTradeResponse
	ctx     context.Context
}

func (a *ExecuteFuturesPerpetualsTradeFuture) Response() (*ExecuteFuturesPerpetualsTradeResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "execute_futures_perpetuals_trade", err)
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

func (r *ExecuteFuturesPerpetualsTradeRequest) Send(ctx context.Context) *ExecuteFuturesPerpetualsTradeFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ExecuteFuturesPerpetualsTradeRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ExecuteFuturesPerpetualsTradeFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ExecuteFuturesPerpetualsTradeResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-binance:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_binance_connection_failed", nil)
		return &ExecuteFuturesPerpetualsTradeFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewBinanceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ExecuteFuturesPerpetualsTrade(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_execute_futures_perpetuals_trade", nil)
			return
		}
		resultc <- rsp
	}()

	return &ExecuteFuturesPerpetualsTradeFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- VerifyCredentials --- //

type VerifyCredentialsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *VerifyCredentialsResponse
	ctx     context.Context
}

func (a *VerifyCredentialsFuture) Response() (*VerifyCredentialsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "verify_credentials", err)
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

func (r *VerifyCredentialsRequest) Send(ctx context.Context) *VerifyCredentialsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *VerifyCredentialsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *VerifyCredentialsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *VerifyCredentialsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-binance:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_binance_connection_failed", nil)
		return &VerifyCredentialsFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewBinanceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.VerifyCredentials(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_verify_credentials", nil)
			return
		}
		resultc <- rsp
	}()

	return &VerifyCredentialsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- List All Asset Pairs --- //

type ListAllAssetPairsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ListAllAssetPairsResponse
	ctx     context.Context
}

func (a *ListAllAssetPairsFuture) Response() (*ListAllAssetPairsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "list_all_asset_pairs", err)
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

func (r *ListAllAssetPairsRequest) Send(ctx context.Context) *ListAllAssetPairsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ListAllAssetPairsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ListAllAssetPairsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ListAllAssetPairsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-binance:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_binance_connection_failed", nil)
		return &ListAllAssetPairsFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewBinanceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ListAllAssetPairs(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "list_all_asset_pairs", nil)
			return
		}
		resultc <- rsp
	}()

	return &ListAllAssetPairsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
