package pager

import "context"

// Pager defines the implementation of pager; which is defined as some mechanism
// of letting some entity know via some event.
type Pager interface {
	// Page pages some entity identified by the identifier.
	// The assumption it is idempotent & rate limited in some way.
	Page(ctx context.Context, identifier, msg string) error
}
