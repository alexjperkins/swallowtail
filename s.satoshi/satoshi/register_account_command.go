package satoshi

import (
	"context"
	"fmt"
	"strings"
	"time"

	accountproto "swallowtail/s.account/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"google.golang.org/grpc"
)

const (
	registerAccountCommandID     = "register-account-command"
	registerAccountCommandPrefix = "!register"
)

func init() {
	registerSatoshiCommand(registerAccountCommandID, handleRegisterAccountCommand)
}

func handleRegisterAccountCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, registerAccountCommandPrefix) {
		return
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	tokens := strings.Split(m.Content, " ")
	if len(tokens) != 3 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: Hi, incorrect usage.\nExample: `!register <username> <password>`"))
		return
	}

	slog.Debug(ctx, "Received %s command, args: %v", registerAccountCommandID, tokens)

	// TODO: validation
	username, password := tokens[1], tokens[2]

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		slog.Error(ctx, "Failed to reach s_account grpc: %v", err)
		return
	}
	defer conn.Close()

	client := accountproto.NewAccountClient(conn)
	if _, err := (client.CreateAccount(ctx, &accountproto.CreateAccountRequest{
		UserId:   m.Author.ID,
		Username: username,
		Password: password,
	})); err != nil {
		slog.Error(ctx, "Failed to create new account: %v", err, map[string]string{
			"user_id":  m.Author.ID,
			"username": username,
		})
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":disappointed: Sorry, I failed to create an account with username: `%s`, maybe you already created one?", username),
		)
		return
	}

	slog.Info(ctx, "Created new account: %s: %s", m.Author.Username, m.Author.ID)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: I have registered your account with username: `%s`", username))
}
