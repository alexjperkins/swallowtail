package main

import (
	"context"
	"fmt"
	"swallowtail/libraries/multiplexing"
	"swallowtail/s.binance/clients"
	"swallowtail/s.binance/consumers"
	"time"

	"github.com/monzo/slog"
)

func main() {
	ctx := context.Background()
	ctx, cf := context.WithCancel(ctx)

	clientMapping := map[string]*clients.StreamClient{}
	consumerGroups := consumers.GetAllBinanceConsumers()

	for endpoint, consumerGroup := range consumerGroups {
		// Client starts the multiplex
		multiplex := multiplexing.New(consumerGroup)
		// Factory class here starts the client
		cli := clients.NewStreamingClient(ctx, endpoint, []*multiplexing.Multiplex{multiplex})
		// Now we have the mapping here
		clientMapping[endpoint] = cli
	}

	t := time.NewTicker(time.Minute * 20)
L:
	for {
		select {
		case <-t.C:
			fmt.Println("Finishing...")
			// Cancel Clients
			for _, c := range clientMapping {
				c.Stop()
			}
			// Cancel Context
			slog.Info(nil, "Cancelling context.")
			cf()
			break L
		}
	}

	time.Sleep(time.Second * 5)
	// TODO: wait group for graceful exit?
}
