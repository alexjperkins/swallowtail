package owner

import (
	"context"
	"fmt"
	"hash/fnv"
	"swallowtail/libraries/ttlcache"
	"time"

	"github.com/monzo/slog"
)

var (
	contactTypeEmail   = "email"
	contactTypeDiscord = "discord"

	// All owners that are defined
	Owners = []*GooglesheetOwner{}
)

type pager func(ctx context.Context, mc MessageClient, msg string) func(context.Context, string) error

type MessageClient interface {
	Send(ctx context.Context, id, msg string) error
	SendPrivateMessage(ctx context.Context, id, msg string) error
}

func New(spreadsheetID string, name string, discordID string, sheetIDs []string, mc MessageClient, minimumTimeBetweenPagers time.Duration, isPrivate bool) *GooglesheetOwner {
	ttl := ttlcache.New(minimumTimeBetweenPagers)
	pager := func(ctx context.Context, msg string) error {
		ttlKey := fmt.Sprintf("%s-%s-%s", name, discordID, hashMsg(msg))
		slog.Error(ctx, "key: %v, ttlcache: %v", ttlKey, ttl)
		if ttl.Exists(ttlKey) {
			return nil
		}
		slog.Info(ctx, "Paging: %s: %s", name, discordID)
		var err error
		if isPrivate {
			err = mc.SendPrivateMessage(ctx, discordID, msg)
		} else {
			err = mc.Send(ctx, discordID, msg)
		}
		if err != nil {
			return err
		}
		// Only set the key if we manage to send
		ttl.SetNull(ttlKey)
		return nil
	}

	return &GooglesheetOwner{
		SpreadsheetID: spreadsheetID,
		SheetsID:      sheetIDs,
		Name:          name,
		DiscordID:     discordID,
		Page:          pager,
		IsPrivate:     isPrivate,
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

func hashMsg(h string) string {
	a := fnv.New32a()
	a.Write([]byte(h))
	return fmt.Sprintf("%v", a.Sum32())
}
