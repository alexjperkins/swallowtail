package gerrors

import (
	"testing"

	"google.golang.org/grpc/codes"
)

// AssertIs is a helper function for asserting gerrors are equal.
func AssertIs(t *testing.T, err error, expectedCode codes.Code, msgs ...string) {
	if !Is(err, expectedCode, msgs...) {
		t.Fatalf("Expected: %v, got: %v", expectedCode, err.Error())
	}
}
