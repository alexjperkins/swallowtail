package commands

import (
	"context"
	"swallowtail/libraries/gerrors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		command     *Command
		withErrCode codes.Code
		withErr     string
	}{
		{
			name: "parent_command_single_arg",
			command: &Command{
				MinimumNumberOfArgs: 0,
				Handler: func(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
					return nil
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.command.Exec(&discordgo.Session{}, &discordgo.MessageCreate{})
			switch {
			case tt.withErr != "":
				require.Error(t, err)

				gerrors.AssertIs(t, err, tt.withErrCode, tt.withErr)
			default:
				require.NoError(t, err)
			}
		})
	}
}
