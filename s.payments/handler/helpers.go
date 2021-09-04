package handler

import (
	"context"
	"fmt"
	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	"time"
)

// CurrentMonthStartTimestamp returns the timestamp of the start of the current month.
// This is defined as the 1st of every month at 00:00:00
func currentMonthStartFromTimestamp(now time.Time) time.Time {
	ts := now.UTC().Truncate(time.Hour)
	daysIntoMonth := ts.Day()

	firstDayOfTheMonth := ts.AddDate(0, 0, -daysIntoMonth+1)
	firstSecondOfTheMonth := firstDayOfTheMonth.Add(-1 * time.Duration(ts.Hour()) * time.Hour)

	return firstSecondOfTheMonth
}

func offboardSubscriber(ctx context.Context, userID, username string) error {
	errParams := map[string]string{
		"user_id": userID,
	}

	// Remove the user as a futures members.
	if err := removeUserAsFuturesMember(ctx, userID); err != nil {
		return gerrors.Augment(err, "failed_to_offboard_user", errParams)
	}

	now := time.Now().UTC()

	// Let them know.
	if _, err := (&discordproto.SendMsgToPrivateChannelRequest{
		UserId:         userID,
		SenderId:       "system:payments",
		Content:        fmt.Sprintf(":disappointed: `Futures Subscription Exipired`.\n Sorry <@%s>, looks like a payment wasn't received for a futures subscription in time.\nPlease ping `@ajperkins` if this is incorrect.", userID),
		IdempotencyKey: fmt.Sprintf("offboardsubscriber-personal-%d-%d", now.Month(), now.Year()),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_offboard_user.notify_user", errParams)
	}

	header := ":rotating_light:   `FUTURES SUB EXPIRED`    :rotating_light:"
	content := `
UserID: %s
Username: %s
Timestamp: %v
	`
	formattedContent := fmt.Sprintf(content, userID, username, time.Now().UTC().Truncate(time.Second))

	// Let us know.
	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiPaymentsPulseChannel,
		SenderId:       "system:payments",
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("offboardsubscriber-pulse-%d-%d", now.Month(), now.Year()),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_offboard_user.publish_to_pulse_channel", errParams)
	}

	return nil
}
