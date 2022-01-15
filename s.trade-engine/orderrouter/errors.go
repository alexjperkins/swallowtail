package orderrouter

import (
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
)

func santizeError(err error) string {
	errorDetail, ok := gerrors.CollectDetailByKeyFromError(err, transport.RequestErrorMessageDetailKey)
	if !ok {
		return "Internal server error: please contact the support channel for help"
	}

	return strings.Join(errorDetail, ": ")
}
