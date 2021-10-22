package ftxproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- List Account Deposits --- //
type ListAccountDepositsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ListAccountDepositsResponse
	ctx     context.Context
}

func (a *ListAccountDepositsFuture) Response() (*ListAccountDepositsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "list_account_deposits", err)
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

func (r *ListAccountDepositsRequest) Send(ctx context.Context) *ListAccountDepositsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ListAccountDepositsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ListAccountDepositsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ListAccountDepositsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-ftx:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_ftx_connection_failed", nil)
		return &ListAccountDepositsFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewFtxClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ListAccountDeposits(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_list_account_deposits", nil)
			return
		}
		resultc <- rsp
	}()

	return &ListAccountDepositsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Get FTX Funding Rates --- //

type GetFTXFundingRatesFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *GetFTXFundingRatesResponse
	ctx     context.Context
}

func (a *GetFTXFundingRatesFuture) Response() (*GetFTXFundingRatesResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_ftx_funding_rate", err)
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

func (r *GetFTXFundingRatesRequest) Send(ctx context.Context) *GetFTXFundingRatesFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetFTXFundingRatesRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetFTXFundingRatesFuture {
	errc := make(chan error, 1)
	resultc := make(chan *GetFTXFundingRatesResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-ftx:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_ftx_connection_failed", nil)
		return &GetFTXFundingRatesFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewFtxClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetFTXFundingRates(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_get_ftx_funding_rate", nil)
			return
		}
		resultc <- rsp
	}()

	return &GetFTXFundingRatesFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
