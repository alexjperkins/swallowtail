package util

import (
	"context"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-multierror"

	"swallowtail/libraries/gerrors"
)

// Retry attempts to execute a function; if it fails, we attempt again until the max number of
// attempts have been acheived. Between each attempt we sleep for some time as determined by a
// exponential backoff algo.
func Retry(ctx context.Context, maxAttempts int, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	var (
		merr error
		rsp  interface{}
		boff = backoff.NewExponentialBackOff()
	)
	for i := 0; i < maxAttempts; i++ {
		r, err := f(ctx)
		if err != nil {
			merr = multierror.Append(merr, err)
			d := boff.NextBackOff()
			time.Sleep(d)
			continue
		}
		rsp = r
	}

	if merr != nil {
		return nil, gerrors.Augment(merr, "failed_to_after_retries_to_execute", map[string]string{
			"attempts": strconv.Itoa(maxAttempts),
		})
	}
	if rsp == nil {
		return nil, gerrors.FailedPrecondition("empty_response", nil)
	}

	return rsp, nil
}
