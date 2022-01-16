package gerrors

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIs_Basic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		err         error
		code        codes.Code
		msgs        []string
		shouldMatch bool
	}{
		{
			name:        "matching_code_with_msg",
			err:         New(ErrAlreadyExists, "account_missing", nil),
			code:        ErrAlreadyExists,
			msgs:        []string{"account_missing"},
			shouldMatch: true,
		},
		{
			name:        "non_matching_code_with_msg",
			err:         New(ErrAlreadyExists, "account_missing", nil),
			code:        ErrCanceled,
			msgs:        []string{"account_missing"},
			shouldMatch: false,
		},
		{
			name:        "matching_code_with_non_matching_msg",
			err:         New(ErrAlreadyExists, "account_missing", nil),
			code:        ErrAlreadyExists,
			msgs:        []string{"bad_message"},
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			didMatch := Is(tt.err, tt.code, tt.msgs...)
			assert.Equal(t, tt.shouldMatch, didMatch)
		})
	}
}

func TestIs_Augmented(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		err          error
		code         codes.Code
		msg          string
		augmentedMsg string
		shouldMatch  bool
	}{
		{
			name:         "matching_code_with_msg",
			err:          New(ErrAlreadyExists, "account_missing", nil),
			msg:          "account_missing",
			augmentedMsg: "failed_to_read_account",
			code:         ErrAlreadyExists,
			shouldMatch:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			aErr := Augment(tt.err, tt.augmentedMsg, nil)

			didMatch := Is(aErr, tt.code, tt.msg, tt.augmentedMsg)
			assert.Equal(t, tt.shouldMatch, didMatch)
		})
	}
}

func TestCollectDetailByKeyFromError(t *testing.T) {
	t.Parallel()

	var (
		userIDKey   = "user_id"
		userIDValue = "harrypotter"
	)

	err := New(ErrBadParam, "missing_param", map[string]string{
		userIDKey: userIDValue,
	})

	details, ok := CollectDetailByKeyFromError(err, userIDKey)
	require.True(t, ok)

	assert.Equal(t, []string{userIDValue}, details)
}

func TestCollectDetailByKeyFromError_Augmented(t *testing.T) {
	t.Parallel()

	var (
		userIDKey = "user_id"

		userIDValue          = "harrypotter"
		userIDValueAugmented = "dumbledore"
	)

	err := New(codes.FailedPrecondition, "missing_param", map[string]string{
		userIDKey: userIDValue,
	})

	err = Augment(err, "failed_request", map[string]string{
		userIDKey: userIDValueAugmented,
	})

	s, ok := status.FromError(err)

	// Basic assertions on the error code in question.
	require.True(t, ok)
	require.Equal(t, s.Code(), codes.FailedPrecondition)

	details, ok := CollectDetailByKeyFromError(err, userIDKey)
	require.True(t, ok)

	sort.Strings(details)
	assert.Equal(t, []string{userIDValueAugmented, userIDValue}, details)
}
