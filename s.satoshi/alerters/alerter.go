package alerters

import "context"

type Alerter interface {
	Start(ctx context.Context)
	Done()
}
