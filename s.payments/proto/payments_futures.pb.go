package paymentsproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- List Account Deposits --- //
type RegisterPaymentFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *RegisterPaymentResponse
	ctx     context.Context
}

func (a *RegisterPaymentFuture) Response() (*RegisterPaymentResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "register_payment", err)
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

func (r *RegisterPaymentRequest) Send(ctx context.Context) *RegisterPaymentFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *RegisterPaymentRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *RegisterPaymentFuture {
	errc := make(chan error, 1)
	resultc := make(chan *RegisterPaymentResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-payments:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_payments_connection_failed", nil)
		return &RegisterPaymentFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewPaymentsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.RegisterPayment(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_register_payment", nil)
			return
		}
		resultc <- rsp
	}()

	return &RegisterPaymentFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
