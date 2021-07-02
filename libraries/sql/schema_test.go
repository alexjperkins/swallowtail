package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseDir(t *testing.T) {
	t.Parallel()
	_, err := baseDir()
	require.NoError(t, err)
}
