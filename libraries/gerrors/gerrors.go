package gerrors

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/monzo/terrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gerrorsproto "swallowtail/libraries/gerrors/proto"
)

// Propagate converts a basic error into a gerror.
func Propagate(code codes.Code, err error, params map[string]string) error {
	return New(code, err.Error(), params)
}

// Is compares gerrors
func Is(err error, code codes.Code, msgs ...string) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	if s.Code() != code {
		return false
	}

	for _, msg := range msgs {
		if !strings.Contains(s.Message(), msg) {
			return false
		}
	}

	return true
}

// Augment augments the given error with a message & extra metadata via params.
func Augment(err error, msg string, params map[string]string) error {
	// Convert error to status.Status
	s, ok := status.FromError(err)
	if !ok {
		return terrors.Augment(err, "Failed to augment gerror", nil)
	}

	details := map[string]string{}
	for k, v := range params {
		details[k] = v
	}

	// Create new status; append message, with the same code & add new metadata details.
	ns := status.Newf(s.Code(), "%s: %s", msg, s.Message())
	ns, e := ns.WithDetails(&gerrorsproto.GerrorMessage{
		Params: details,
	})
	if e != nil {
		return e
	}

	// Update with old metadata.
	for _, v := range s.Details() {
		m, ok := v.(proto.Message)
		if !ok {
			// Best effort
			continue
		}

		dv, e := ptypes.MarshalAny(m)
		if e != nil {
			return e
		}

		ns, e = ns.WithDetails(dv)
		if e != nil {
			return e
		}
	}

	return ns.Err()
}

// New ...
func New(code codes.Code, msg string, params map[string]string) error {
	s := status.New(code, msg)

	s, err := s.WithDetails(&gerrorsproto.GerrorMessage{
		Params: params,
	})
	if err != nil {
		return err
	}

	return s.Err()
}

// NotFound ...
func NotFound(msg string, params map[string]string) error {
	return New(ErrNotFound, msg, params)
}

// AlreadyExists ...
func AlreadyExists(msg string, params map[string]string) error {
	return New(ErrAlreadyExists, msg, params)
}

// FailedPrecondition ...
func FailedPrecondition(msg string, params map[string]string) error {
	return New(ErrFailedPrecondition, msg, params)
}

// BadParam ...
func BadParam(msg string, params map[string]string) error {
	return New(ErrBadParam, msg, params)
}

// Unauthenticated ...
func Unauthenticated(msg string, params map[string]string) error {
	return New(ErrUnauthenticated, msg, params)
}
