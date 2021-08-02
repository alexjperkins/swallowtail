package sync

import (
	"context"

	"github.com/hashicorp/go-multierror"
)

// Syncer defines the contract for sync a googlesheets
type Syncer interface {
	// Sync sync all registered users sheets
	Sync(context.Context) error
	// Reload reloads all stored list of users sheets if required.
	Refresh(context.Context) error
}

func Init(ctx context.Context) error {
	mu.RLock()
	defer mu.RLock()
	var err error
	for _, syncer := range registry {
		e := syncer.Refresh(ctx)
		if e != nil {
			multierror.Append(err, e)
		}
		go syncer.Sync(ctx)
	}

	return err
}
