package main

import (
	"context"
	"fmt"
	"time"

	"github.com/monzo/slog"
	"google.golang.org/grpc"

	coingeckoproto "swallowtail/s.coingecko/proto"
)

func main() {
	rsp, err := GetAssetLatestPriceByIDRequest{
		AssetID:   "btc",
		AssetPair: "usdt",
	}.SendWithTimeout(context.Background(), 30*time.Second).Response()
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.LatestPrice)
}

type GetAssetLatestPriceByIDRequest struct {
	AssetID   string
	AssetPair string
}

type GetAssetLatestPriceByIDFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *coingeckoproto.GetAssetLatestPriceByIDResponse
	ctx     context.Context
}

func (a *GetAssetLatestPriceByIDFuture) Response() (*coingeckoproto.GetAssetLatestPriceByIDResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "get_latest_asset_price_by_id", err)
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

func (r GetAssetLatestPriceByIDRequest) Send(ctx context.Context) *GetAssetLatestPriceByIDFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r GetAssetLatestPriceByIDRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *GetAssetLatestPriceByIDFuture {
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

		rsp, err := c.GetAssetLatestPriceByID(ctx, &coingeckoproto.GetAssetLatestPriceByIDRequest{
			CoingeckoCoinId: r.AssetID,
			AssetPair:       r.AssetPair,
		})
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
