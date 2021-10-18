package marketdataproto

import (
	context "context"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Publish Latest Price Information --- //

type PublishLatestPriceInformationFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *PublishLatestPriceInformationResponse
	ctx     context.Context
}

func (a *PublishLatestPriceInformationFuture) Response() (*PublishLatestPriceInformationResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "publish_latest_price_information", err)
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

func (r *PublishLatestPriceInformationRequest) Send(ctx context.Context) *PublishLatestPriceInformationFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *PublishLatestPriceInformationRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *PublishLatestPriceInformationFuture {
	errc := make(chan error, 1)
	resultc := make(chan *PublishLatestPriceInformationResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-marketdata:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &PublishLatestPriceInformationFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewMarketdataClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.PublishLatestPriceInformation(ctx, r)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &PublishLatestPriceInformationFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
