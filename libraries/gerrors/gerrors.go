package gerrors

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gerrorsproto "swallowtail/libraries/gerrors/proto"
)

// Propagate converts a basic error into a gerror.
func Propagate(err error, code codes.Code, params map[string]string) error {
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
		return status.Error(ErrUnknown, "Failed to augment gerror")
	}

	var e error
	s, e = s.WithDetails(&gerrorsproto.GerrorMessage{
		Params: params,
	})
	if e != nil {
		return e
	}

	return s.Err()
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

// CollectDetailByKeyFromError ...
func CollectDetailByKeyFromError(err error, key string) ([]string, bool) {
	st, ok := status.FromError(err)
	if !ok {
		return nil, false
	}

	var ss []string
	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *gerrorsproto.GerrorMessage:
			s, ok := d.Params[key]
			if !ok {
				continue
			}

			ss = append(ss, s)
		}
	}

	return ss, true
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

func Unimplemented(msg string, params map[string]string) error {
	return New(ErrUnimplemented, msg, params)
}
