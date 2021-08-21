package accountproto

import (
	context "context"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

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
			errc <- err
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

// --- Add Exchange --- //

type AddExchangeFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *AddExchangeResponse
	ctx     context.Context
}

func (a *AddExchangeFuture) Response() (*AddExchangeResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "add_exchange", err)
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

func (r *AddExchangeRequest) Send(ctx context.Context) *AddExchangeFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *AddExchangeRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *AddExchangeFuture {
	errc := make(chan error, 1)
	resultc := make(chan *AddExchangeResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &AddExchangeFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.AddExchange(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &AddExchangeFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

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
