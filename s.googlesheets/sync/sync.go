package sync

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"
)

// Syncer defines the contract for sync a googlesheets
type Syncer interface {
	// Sync sync all registered users sheets
	Sync(context.Context)
	// Reload reloads all stored list of users sheets if required.
	Refresh(context.Context) error
}

func Init(ctx context.Context) error {
	var err error
	for id, syncer := range registry {
		e := syncer.Refresh(ctx)
		if e != nil {
			multierror.Append(err, e)
		}

		slog.Debug(ctx, "Syncer: %s initialized", id)
		go syncer.Sync(ctx)
	}

	return err
}
