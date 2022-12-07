package cassandra

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateKeyspace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		keyspaceName string

		expectedError error
		requireF      func(t require.TestingT, err error, msgAndArgs ...interface{})
	}{
		{
			name:          "invalid_keyspace_empty",
			keyspaceName:  "",
			expectedError: ErrKeyspaceNameEmpty,
			requireF:      require.Error,
		},
		{
			name:         "valid_keyspace_name_under_boundary",
			keyspaceName: strings.Repeat("a", maximumKeyspaceLength-1),
			requireF:     require.NoError,
		},
		{
			name:         "valid_keyspace_name_on_boundary",
			keyspaceName: strings.Repeat("a", maximumKeyspaceLength),
			requireF:     require.NoError,
		},
		{
			name:          "invalid_keyspace_name_too_long",
			keyspaceName:  strings.Repeat("a", maximumKeyspaceLength+1),
			expectedError: ErrKeyspaceNameTooLong,
			requireF:      require.Error,
		},
		{
			name:          "invalid_keyspace_name_invalid_char",
			keyspaceName:  "a_keyspace_invalid_?",
			expectedError: ErrKeyspaceNameInvalidChar,
			requireF:      require.Error,
		},
		{
			name:         "valid_keyspace_name",
			keyspaceName: "a_keyspace_valid",
			requireF:     require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateKeyspace(tt.keyspaceName)
			tt.requireF(t, err)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
