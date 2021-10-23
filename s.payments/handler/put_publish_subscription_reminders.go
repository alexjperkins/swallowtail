package handler

import (
	"context"
	"fmt"
	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	paymentsproto "swallowtail/s.payments/proto"
	"time"
)

// PublishSubscriptionReminder ...
func (s *PaymentsService) PublishSubscriptionReminder(
	ctx context.Context, in *paymentsproto.PublishSubscriptionReminderRequest,
) (*paymentsproto.PublishSubscriptionReminderResponse, error) {
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

	errParams := map[string]string{
		"actor_id":          in.ActorId,
		"subscription_type": in.ReminderType.String(),
	}

	// Validate the caller is authorized to call this RPC.
	validActor, err := isActorValid(ctx, in.ActorId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions.actor_check", errParams)
	}
	if !validActor {
		return nil, gerrors.Unauthenticated("failed_to_enforce_subscriptions.unauthorized", errParams)
	}

	var reminder string
	switch in.ReminderType {
	case paymentsproto.SubscriptionReminderType_MINUS_54_HOURS:
		reminder = "tomorrow before 5 minutes before midnight"
	case paymentsproto.SubscriptionReminderType_MINUS_4_HOURS:
		reminder = "in 4 hours"
	case paymentsproto.SubscriptionReminderType_MINUS_1_HOUR:
		reminder = "in 1 hour"
	default:
		return nil, gerrors.FailedPrecondition("failed_to_publish_subscription_reminders.invalid_type", errParams)
	}

	// Idempotent on the month, year & the subscription message type.
	now := time.Now().UTC()
	idempotencyKey := fmt.Sprintf("spaymentsreminder-%s-%d-%d", in.ReminderType.String(), now.Month(), now.Year())

	if _, err := (&discordproto.SendMsgToChannelRequest{
		Content:        formatReminderMsg(reminder),
		SenderId:       "system:payments",
		ChannelId:      discordproto.DiscordSatoshiPaymentsPulseChannel,
		IdempotencyKey: idempotencyKey,
		Force:          in.Force,
	}).Send(ctx).Response(); err != nil {
		return nil, gerrors.Augment(err, "failed_to_publish_subscription_reminders", errParams)
	}

	// We ideally need to return success; but we don't know if it actually sent since we dont' return this from `s.discord`
	// we should update.
	return &paymentsproto.PublishSubscriptionReminderResponse{}, nil
}

func formatReminderMsg(reminder string) string {
	return fmt.Sprintf(
		":wave:   Headsup @everyone   :rotating_light:\nFriendly reminder that futures subscriptions are due **%s**.\nPlease make sure you have registered a payment before then. Thanks :pray:",
		reminder,
	)
}
