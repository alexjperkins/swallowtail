package binanceproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Execute New Futures Perpetual Order --- //

type ExecuteNewFuturesPerpetualOrderFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ExecuteNewFuturesPerpetualOrderResponse
	ctx     context.Context
}

func (a *ExecuteNewFuturesPerpetualOrderFuture) Response() (*ExecuteNewFuturesPerpetualOrderResponse, error) {
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

func (r *ExecuteNewFuturesPerpetualOrderRequest) Send(ctx context.Context) *ExecuteNewFuturesPerpetualOrderFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ExecuteNewFuturesPerpetualOrderRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ExecuteNewFuturesPerpetualOrderFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ExecuteNewFuturesPerpetualOrderResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-binance:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_binance_connection_failed", nil)
		return &ExecuteNewFuturesPerpetualOrderFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewBinanceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ExecuteNewFuturesPerpetualOrder(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_execute_futures_perpetuals_trade", nil)
			return
		}
		resultc <- rsp
	}()

	return &ExecuteNewFuturesPerpetualOrderFuture{
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

// --- Get Latest Price --- //

type GetLatestPricesFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *GetLatestPriceResponse
	ctx     context.Context
}

func (a *GetLatestPricesFuture) Response() (*GetLatestPriceResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_latest_price", err)
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

func (r *GetLatestPriceRequest) Send(ctx context.Context) *GetLatestPricesFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetLatestPriceRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetLatestPricesFuture {
	errc := make(chan error, 1)
	resultc := make(chan *GetLatestPriceResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-binance:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_binance_connection_failed", nil)
		return &GetLatestPricesFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewBinanceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetLatestPrice(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "get_latest_prices", nil)
			return
		}
		resultc <- rsp
	}()

	return &GetLatestPricesFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Read Perpetual Futures Account --- //

type ReadPerpetualFuturesAccountsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ReadPerpetualFuturesAccountResponse
	ctx     context.Context
}

func (a *ReadPerpetualFuturesAccountsFuture) Response() (*ReadPerpetualFuturesAccountResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_perpetual_futures_account", err)
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

func (r *ReadPerpetualFuturesAccountRequest) Send(ctx context.Context) *ReadPerpetualFuturesAccountsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ReadPerpetualFuturesAccountRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ReadPerpetualFuturesAccountsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ReadPerpetualFuturesAccountResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-binance:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_binance_connection_failed", nil)
		return &ReadPerpetualFuturesAccountsFuture{
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
	c := NewBinanceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ReadPerpetualFuturesAccount(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_read_perpetual_futures_account", nil)
			return
		}
		resultc <- rsp
	}()

	return &ReadPerpetualFuturesAccountsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Read Funding Rates --- //

type GetFundingRatesFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *GetFundingRatesResponse
	ctx     context.Context
}

func (a *GetFundingRatesFuture) Response() (*GetFundingRatesResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_funding_rate", err)
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

func (r *GetFundingRatesRequest) Send(ctx context.Context) *GetFundingRatesFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetFundingRatesRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetFundingRatesFuture {
	errc := make(chan error, 1)
	resultc := make(chan *GetFundingRatesResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-binance:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_binance_connection_failed", nil)
		return &GetFundingRatesFuture{
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
	c := NewBinanceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetFundingRates(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_get_funding_rate", nil)
			return
		}
		resultc <- rsp
	}()

	return &GetFundingRatesFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
