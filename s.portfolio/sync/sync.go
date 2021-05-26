package sync

import "context"

// Syncer defines the contract for sync a googlesheets
type Syncer interface {
	Start(context.Context)
	Done()
}
