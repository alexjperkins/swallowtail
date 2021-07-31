package main

import (
	"context"
	"time"

	"google.golang.org/grpc"

	coingeckoproto "swallowtail/s.coingecko/proto"
)

type GetAssetLatestPriceByIDFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *coingeckoproto.GetAssetLatestPriceByIDResponse
	ctx     context.Context
}

func (a *GetAssetLatestPriceByIDFuture) Response() (*coingeckoproto.GetAssetLatestPriceByIDResponse, error) {
	// What if this errors?
	defer a.closer()

	select {
	case r := <-a.resultc:
		return r, nil
	case <-a.ctx.Done():
		return nil, a.ctx.Err()
	case err := <-a.errc:
		return nil, err
	}
}

func Send(ctx context.Context, in *coingeckoproto.GetAssetLatestPriceByIDRequest) *GetAssetLatestPriceByIDFuture {
	return SendWithTimeout(ctx, in, 10*time.Second)
}

func SendWithTimeout(ctx context.Context, in *coingeckoproto.GetAssetLatestPriceByIDRequest, timeout time.Duration) *GetAssetLatestPriceByIDFuture {
	errc := make(chan error, 1)
	resultc := make(chan *coingeckoproto.GetAssetLatestPriceByIDResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-coingecko:8000", grpc.WithInsecure())
	if err != nil {
		errc <- err
		return &GetAssetLatestPriceByIDFuture{
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := coingeckoproto.NewCoingeckoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		defer func() {
			cancel()
		}()

		rsp, err := c.GetAssetLatestPriceByID(ctx, in)
		if err != nil {
			errc <- err
			return
		}
		resultc <- rsp
	}()

	return &GetAssetLatestPriceByIDFuture{
		closer: conn.Close,
		errc:   errc,
		ctx:    ctx,
	}
}
