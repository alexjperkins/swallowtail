package ftxproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Get FTX Status --- //
type GetFTXStatus struct {
	closer  func() error
	errc    chan error
	resultc chan *GetFTXStatusResponse
	ctx     context.Context
}

func (a *GetFTXStatus) Response() (*GetFTXStatusResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_ftx_status", err)
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

func (r *GetFTXStatusRequest) Send(ctx context.Context) *GetFTXStatus {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetFTXStatusRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetFTXStatus {
	errc := make(chan error, 1)
	resultc := make(chan *GetFTXStatusResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-ftx:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s-ftx_connection_failed", nil)
		return &GetFTXStatus{
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
	c := NewFtxClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetFTXStatus(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_get_ftx_status", nil)
			return
		}
		resultc <- rsp
	}()

	return &GetFTXStatus{
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
type GetFTXFundingRates struct {
	closer  func() error
	errc    chan error
	resultc chan *GetFTXFundingRatesResponse
	ctx     context.Context
}

func (a *GetFTXFundingRates) Response() (*GetFTXFundingRatesResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_ftx_funding_rates", err)
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

func (r *GetFTXFundingRatesRequest) Send(ctx context.Context) *GetFTXFundingRates {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *GetFTXFundingRatesRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetFTXFundingRates {
	errc := make(chan error, 1)
	resultc := make(chan *GetFTXFundingRatesResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-ftx:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s-ftx_connection_failed", nil)
		return &GetFTXFundingRates{
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
	c := NewFtxClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.GetFTXFundingRates(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_get_ftx_funding_rates", nil)
			return
		}
		resultc <- rsp
	}()

	return &GetFTXFundingRates{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- List Account Deposits --- //

type ListAccountDeposits struct {
	closer  func() error
	errc    chan error
	resultc chan *ListAccountDepositsResponse
	ctx     context.Context
}

func (a *ListAccountDeposits) Response() (*ListAccountDepositsResponse, error) {
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

func (r *ListAccountDepositsRequest) Send(ctx context.Context) *ListAccountDeposits {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ListAccountDepositsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ListAccountDeposits {
	errc := make(chan error, 1)
	resultc := make(chan *ListAccountDepositsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-ftx:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s-ftx_connection_failed", nil)
		return &ListAccountDeposits{
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

	return &ListAccountDeposits{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Execute New Order --- //

type ExecuteNewOrder struct {
	closer  func() error
	errc    chan error
	resultc chan *ExecuteNewOrderResponse
	ctx     context.Context
}

func (a *ExecuteNewOrder) Response() (*ExecuteNewOrderResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "execute_new_order", err)
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

func (r *ExecuteNewOrderRequest) Send(ctx context.Context) *ExecuteNewOrder {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ExecuteNewOrderRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ExecuteNewOrder {
	errc := make(chan error, 1)
	resultc := make(chan *ExecuteNewOrderResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-ftx:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s-ftx_connection_failed", nil)
		return &ExecuteNewOrder{
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
	c := NewFtxClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ExecuteNewOrder(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_execute_new_order", nil)
			return
		}
		resultc <- rsp
	}()

	return &ExecuteNewOrder{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- List FTX Instruments --- //

type ListFTXInstruments struct {
	closer  func() error
	errc    chan error
	resultc chan *ListFTXInstrumentsResponse
	ctx     context.Context
}

func (a *ListFTXInstruments) Response() (*ListFTXInstrumentsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "list_ftx_instruments", err)
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

func (r *ListFTXInstrumentsRequest) Send(ctx context.Context) *ListFTXInstruments {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ListFTXInstrumentsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ListFTXInstruments {
	errc := make(chan error, 1)
	resultc := make(chan *ListFTXInstrumentsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-ftx:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s-ftx_connection_failed", nil)
		return &ListFTXInstruments{
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
	c := NewFtxClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ListFTXInstruments(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_list_ftx_instruments", nil)
			return
		}
		resultc <- rsp
	}()

	return &ListFTXInstruments{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
