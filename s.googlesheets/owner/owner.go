package owner

import (
	"context"

	"github.com/monzo/slog"
)

var (
	contactTypeEmail   = "email"
	contactTypeDiscord = "discord"
)

type pager func(ctx context.Context, mc MessageClient, msg string) func(context.Context, string) error

type MessageClient interface {
	Send(ctx context.Context, id, msg string) error
	SendPrivateMessage(ctx context.Context, id, msg string) error
}

func New(spreadsheetID string, name string, discordID string, sheetIDs []string, mc MessageClient, isPrivate bool) *GooglesheetOwner {
	return &GooglesheetOwner{
		SpreadsheetID: spreadsheetID,
		SheetsID:      sheetIDs,
		Name:          name,
		DiscordID:     discordID,
		Page: func(ctx context.Context, msg string) error {
			slog.Info(ctx, "Paging: %s: %s", name, discordID)
			if isPrivate {
				return mc.SendPrivateMessage(ctx, discordID, msg)
			}
			return mc.Send(ctx, discordID, msg)
		},
		IsPrivate: isPrivate,
	}
}

type GooglesheetOwner struct {
	SpreadsheetID string
	SheetsID      []string
	// Name of the owner of the googlesheets client
	Name string
	// DiscordID the discord id for the given user.
	DiscordID string
	// Pager pings owner
	Page      func(ctx context.Context, msg string) error
	IsPrivate bool
}
