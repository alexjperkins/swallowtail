package pager

import "context"

type Pager interface {
	Page(ctx context.Context, identifier, msg string) error
}
