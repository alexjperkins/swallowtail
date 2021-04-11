package sync

import "context"

// GoogleSheetsSyncer defines the contract for sync a googlesheets
type GoogleSheetsSyncer interface {
	Start(context.Context)
}
