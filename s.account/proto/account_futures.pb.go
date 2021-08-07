package accountproto

import (
	context "context"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- ReadAccount --- //
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
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		defer func() {
			cancel()
		}()

		rsp, err := c.ReadAccount(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &ReadAccountFuture{
		closer: conn.Close,
		errc:   errc,
		ctx:    ctx,
	}
}
