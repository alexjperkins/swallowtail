package solananftsproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- Read Solana Price Statistics By Collection ID --- //
type ReadSolanaPriceStatisticsByCollectionIDFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ReadSolanaPriceStatisticsByCollectionIDResponse
	ctx     context.Context
}

func (a *ReadSolanaPriceStatisticsByCollectionIDFuture) Response() (*ReadSolanaPriceStatisticsByCollectionIDResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_solana_price_statistics_by_collection_id", err)
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

func (r *ReadSolanaPriceStatisticsByCollectionIDRequest) Send(ctx context.Context) *ReadSolanaPriceStatisticsByCollectionIDFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ReadSolanaPriceStatisticsByCollectionIDRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ReadSolanaPriceStatisticsByCollectionIDFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ReadSolanaPriceStatisticsByCollectionIDResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-solananfts:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_solananfts_connection_failed", nil)
		return &ReadSolanaPriceStatisticsByCollectionIDFuture{
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
	c := NewSolananftsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ReadSolanaPriceStatisticsByCollectionID(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_read_solana_price_statistics_by_collection_id", nil)
			return
		}
		resultc <- rsp
	}()

	return &ReadSolanaPriceStatisticsByCollectionIDFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
