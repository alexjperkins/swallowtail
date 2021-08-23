package gerrors

import "google.golang.org/grpc/codes"

const (
	ErrFailedPrecondition = codes.FailedPrecondition
	ErrNotFound           = codes.NotFound
	ErrUnauthenticated    = codes.Unauthenticated
	ErrAlreadyExists      = codes.AlreadyExists
	ErrUnknown            = codes.Unknown
	ErrUnimplemented      = codes.Unimplemented
	ErrUnavailable        = codes.Unavailable
	ErrPermissionDenied   = codes.PermissionDenied
	ErrCanceled           = codes.Canceled
	ErrDeadlineExceeded   = codes.DeadlineExceeded
	ErrBadParam           = codes.InvalidArgument
	ErrRateLimited        = codes.ResourceExhausted
)
