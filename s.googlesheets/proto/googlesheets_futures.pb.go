package googlesheetsproto

import (
	context "context"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- CreatePortfolioSheet --- //
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

	conn, err := grpc.DialContext(ctx, "swallowtail-s-googlesheets:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &CreatePortfolioSheetFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewGooglesheetsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.CreatePortfolioSheet(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &CreatePortfolioSheetFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- ListSheetByUserID --- //

type ListSheetByUserIDFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ListSheetsByUserIDResponse
	ctx     context.Context
}

func (a *ListSheetByUserIDFuture) Response() (*ListSheetsByUserIDResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "list_sheet_user_id", err)
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

func (r *ListSheetsByUserIDRequest) Send(ctx context.Context) *ListSheetByUserIDFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ListSheetsByUserIDRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ListSheetByUserIDFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ListSheetsByUserIDResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-googlesheets:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ListSheetByUserIDFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewGooglesheetsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ListSheetsByUserID(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &ListSheetByUserIDFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
