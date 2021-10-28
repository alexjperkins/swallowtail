package bitfinexproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Get Bitfinex Funding Rates --- //

type GetBitfinexFundingRatesFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *GetBitfinexFundingRatesResponse
	ctx     context.Context
}

func (a *GetBitfinexFundingRatesFuture) Response() (*GetBitfinexFundingRatesResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_bitfinex_funding_rates", err)
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

func (r *GetBitfinexFundingRatesRequest) Send(ctx context.Context) *GetBitfinexFundingRatesFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetBitfinexFundingRatesRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetBitfinexFundingRatesFuture {
	errc := make(chan error, 1)
	resultc := make(chan *GetBitfinexFundingRatesResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-bitfinex:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_bitfinex_connection_failed", nil)
		return &GetBitfinexFundingRatesFuture{
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
	c := NewBitfinexClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetBitfinexFundingRates(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_get_bitfinex_funding_rates", nil)
			return
		}
		resultc <- rsp
	}()

	return &GetBitfinexFundingRatesFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
