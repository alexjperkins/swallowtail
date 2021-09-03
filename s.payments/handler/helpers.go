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
func currentMonthStartTimestamp() time.Time {
	now := time.Now().UTC().Truncate(time.Hour)
	daysIntoMonth := now.Day()
	return now.AddDate(0, 0, -daysIntoMonth)
}

func offboardSubscriber(ctx context.Context, userID string) error {
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
		Content:        ":disappointed: Sorry <@%s>, looks like a payment wasn't received for futures subscription in time. Please ping @ajperkins if this is incorrect.",
		IdempotencyKey: fmt.Sprintf("offboardsubscriber-%d-%d", now.Month(), now.Year()),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_offboard_user.notify_user", errParams)
	}

	// Let us know.
	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiPaymentsPulseChannel,
		SenderId:       "system:payments",
		Content:        ":rotating_light: Subscriber `%s` hasn't registered a payment for futures subscription. They have been offboard :grimacing:",
		IdempotencyKey: fmt.Sprintf("offboardsubscriber-%d-%d", now.Month(), now.Year()),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_offboard_user.publish_to_pulse_channel", errParams)
	}

	return nil
}
