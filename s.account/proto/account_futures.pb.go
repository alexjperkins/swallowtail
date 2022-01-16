package accountproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Create Account --- //

type CreateAccountFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *CreateAccountResponse
	ctx     context.Context
}

func (a *CreateAccountFuture) Response() (*CreateAccountResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "create_account", err)
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

func (r *CreateAccountRequest) Send(ctx context.Context) *CreateAccountFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *CreateAccountRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *CreateAccountFuture {
	errc := make(chan error, 1)
	resultc := make(chan *CreateAccountResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &CreateAccountFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.CreateAccount(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &CreateAccountFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- List Account --- //

type ListAccountsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ListAccountsResponse
	ctx     context.Context
}

func (a *ListAccountsFuture) Response() (*ListAccountsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "list_accounts", err)
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

func (r *ListAccountsRequest) Send(ctx context.Context) *ListAccountsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ListAccountsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ListAccountsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ListAccountsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ListAccountsFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ListAccounts(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &ListAccountsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Read Account --- //

type ReadAccountFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ReadAccountResponse
	ctx     context.Context
}

func (a *ReadAccountFuture) Response() (*ReadAccountResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_account", err)
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

func (r *ReadAccountRequest) Send(ctx context.Context) *ReadAccountFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ReadAccountRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ReadAccountFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ReadAccountResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ReadAccountFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ReadAccount(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &ReadAccountFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Update Account --- //

type UpdateAccountFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *UpdateAccountResponse
	ctx     context.Context
}

func (a *UpdateAccountFuture) Response() (*UpdateAccountResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "update_account", err)
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

func (r *UpdateAccountRequest) Send(ctx context.Context) *UpdateAccountFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *UpdateAccountRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *UpdateAccountFuture {
	errc := make(chan error, 1)
	resultc := make(chan *UpdateAccountResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &UpdateAccountFuture{
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
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.UpdateAccount(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &UpdateAccountFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Page Account --- //

type PageAccountFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *PageAccountResponse
	ctx     context.Context
}

func (a *PageAccountFuture) Response() (*PageAccountResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "page_account", err)
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

func (r *PageAccountRequest) Send(ctx context.Context) *PageAccountFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *PageAccountRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *PageAccountFuture {
	errc := make(chan error, 1)
	resultc := make(chan *PageAccountResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &PageAccountFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.PageAccount(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_update_account", nil)
			return
		}
		resultc <- rsp
	}()

	return &PageAccountFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Add Venue Account --- //

type AddVenueAccountFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *AddVenueAccountResponse
	ctx     context.Context
}

func (a *AddVenueAccountFuture) Response() (*AddVenueAccountResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "add_venue_account", err)
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

func (r *AddVenueAccountRequest) Send(ctx context.Context) *AddVenueAccountFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *AddVenueAccountRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *AddVenueAccountFuture {
	errc := make(chan error, 1)
	resultc := make(chan *AddVenueAccountResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &AddVenueAccountFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.AddVenueAccount(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &AddVenueAccountFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Read Primary Venue Account By User ID--- //

type ReadPrimaryVenueAccountByUserIDFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ReadPrimaryVenueAccountByUserIDResponse
	ctx     context.Context
}

func (a *ReadPrimaryVenueAccountByUserIDFuture) Response() (*ReadPrimaryVenueAccountByUserIDResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_primary_venue_account_by_user_id", err)
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

func (r *ReadPrimaryVenueAccountByUserIDRequest) Send(ctx context.Context) *ReadPrimaryVenueAccountByUserIDFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ReadPrimaryVenueAccountByUserIDRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ReadPrimaryVenueAccountByUserIDFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ReadPrimaryVenueAccountByUserIDResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ReadPrimaryVenueAccountByUserIDFuture{
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
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ReadPrimaryVenueAccountByUserID(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_read_primary_venue_account_by_user_id", nil)
			return
		}
		resultc <- rsp
	}()

	return &ReadPrimaryVenueAccountByUserIDFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Read Venue Account By VenueAccount Details --- //

type ReadVenueAccountByVenueAccountDetailsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ReadVenueAccountByVenueAccountDetailsResponse
	ctx     context.Context
}

func (a *ReadVenueAccountByVenueAccountDetailsFuture) Response() (*ReadVenueAccountByVenueAccountDetailsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_venue_account_by_venue_account_details", err)
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

func (r *ReadVenueAccountByVenueAccountDetailsRequest) Send(ctx context.Context) *ReadVenueAccountByVenueAccountDetailsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ReadVenueAccountByVenueAccountDetailsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ReadVenueAccountByVenueAccountDetailsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ReadVenueAccountByVenueAccountDetailsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ReadVenueAccountByVenueAccountDetailsFuture{
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
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ReadVenueAccountByVenueAccountDetails(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_read_venue_account_by_venue_account_details", nil)
			return
		}
		resultc <- rsp
	}()

	return &ReadVenueAccountByVenueAccountDetailsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- List Venue Accounts --- //

type ListVenueAccountsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ListVenueAccountsResponse
	ctx     context.Context
}

func (a *ListVenueAccountsFuture) Response() (*ListVenueAccountsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "list_venue_accounts", err)
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

func (r *ListVenueAccountsRequest) Send(ctx context.Context) *ListVenueAccountsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ListVenueAccountsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ListVenueAccountsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ListVenueAccountsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &ListVenueAccountsFuture{
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
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ListVenueAccounts(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_list_venue_accounts", nil)
			return
		}
		resultc <- rsp
	}()

	return &ListVenueAccountsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
