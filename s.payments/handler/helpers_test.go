package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCurrentMonthStartTimestamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		now  time.Time
	}{
		{
			name: "random-now",
			now:  time.Now(),
		},
		{
			name: "last-second-of-the-month",
			now:  time.Date(2020, time.November, 30, 23, 59, 59, 0, time.Local),
		},
		{
			name: "first-second-of-the-month",
			now:  time.Date(2020, time.December, 1, 00, 00, 1, 0, time.Local),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := currentMonthStartFromTimestamp(tt.now)

			assert.Equal(t, 1, res.Day())
			assert.Equal(t, tt.now.Month(), res.Month())

			assert.Equal(t, 0, res.Hour())
			assert.Equal(t, 0, res.Minute())
			assert.Equal(t, 0, res.Second())
		})
	}
}
