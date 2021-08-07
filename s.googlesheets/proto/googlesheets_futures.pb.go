package googlesheetsproto

import (
	context "context"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

///  --- CreatePortfolioSheet --- ///

type CreatePortfolioSheetFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *CreatePortfolioSheetResponse
	ctx     context.Context
}

func (a *CreatePortfolioSheetFuture) Response() (*CreatePortfolioSheetResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "create_portfolio_sheet", err)
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

func (r *CreatePortfolioSheetRequest) Send(ctx context.Context) *CreatePortfolioSheetFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *CreatePortfolioSheetRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *CreatePortfolioSheetFuture {
	errc := make(chan error, 1)
	resultc := make(chan *CreatePortfolioSheetResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-coingecko:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &CreatePortfolioSheetFuture{
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewGooglesheetsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		defer func() {
			cancel()
		}()

		rsp, err := c.CreatePortfolioSheet(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &CreatePortfolioSheetFuture{
		closer: conn.Close,
		errc:   errc,
		ctx:    ctx,
	}
}

///  --- ListSheetsByUserID --- ///

type ListSheetsByUserIDFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ListSheetsByUserIDResponse
	ctx     context.Context
}

func (a *ListSheetsByUserIDFuture) Response() (*ListSheetsByUserIDResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "create_portfolio_sheet", err)
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

func (r *ListSheetsByUserIDRequest) Send(ctx context.Context) *ListSheetsByUserIDFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ListSheetsByUserIDRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ListSheetsByUserIDFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ListSheetsByUserIDResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-coingecko:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ListSheetsByUserIDFuture{
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewGooglesheetsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		defer func() {
			cancel()
		}()

		rsp, err := c.ListSheetsByUserID(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &ListSheetsByUserIDFuture{
		closer: conn.Close,
		errc:   errc,
		ctx:    ctx,
	}
}