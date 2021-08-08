package gerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
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

// TODO: Test params under augment.
