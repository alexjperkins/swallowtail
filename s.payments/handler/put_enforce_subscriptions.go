package handler

import (
	"context"
	"fmt"
	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.payments/dao"
	paymentsproto "swallowtail/s.payments/proto"
	"time"

	"github.com/monzo/slog"
)

// EnforceSubscriptions ...
func (s *PaymentsService) EnforceSubscriptions(
	ctx context.Context, in *paymentsproto.EnforceSubscriptionsRequest,
) (*paymentsproto.EnforceSubscriptionsResponse, error) {
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

	// TODO validate the day.

	errParams := map[string]string{
		"actor_id": in.ActorId,
	}

	// Validate the caller is authorized to call this RPC.
	validActor, err := isActorValid(ctx, in.ActorId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions.actor_check", errParams)
	}
	if !validActor {
		return nil, gerrors.Unauthenticated("failed_to_enforce_subscriptions.unauthorized", errParams)
	}

	// Notify Pulse channel.
	if _, err := (&discordproto.SendMsgToChannelRequest{
		Content:   formatCronitorMsg("Enforce Subscriptions", in.ActorId, "Started", time.Now().Truncate(time.Second)),
		ChannelId: discordproto.DiscordSatoshiPaymentsPulseChannel,
		SenderId:  "system:payments",
		Force:     true,
	}).Send(ctx).Response(); err != nil {
		return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions.failed_to_notify_pulse_channel", errParams)
	}

	futuresMembers, err := listFuturesMembers(ctx)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions", errParams)
	}

	for _, fm := range futuresMembers {
		if fm.IsAdmin {
			slog.Warn(ctx, "Skipping subscription payment check for admin: %s: %s", fm.UserId, fm.Username)
			continue
		}

		errParams["user_id"] = fm.UserId

		ok, err := dao.UserPaymentExistsSince(ctx, fm.UserId, currentMonthStartFromTimestamp(time.Now()))
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions", errParams)
		}

		// Futures member has paid for the month; we can continue.
		if ok {
			continue
		}

		slog.Info(ctx, "Offboarding user: %v: %v", fm.Username, fm.UserId)

		// Uh-oh they haven't paid, lets offboard them.
		if err := offboardSubscriber(ctx, fm.UserId, fm.Username); err != nil {
			return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions", errParams)
		}
	}

	return &paymentsproto.EnforceSubscriptionsResponse{}, nil
}

func isActorValid(ctx context.Context, actorID string) (bool, error) {
	switch actorID {
	case paymentsproto.ActorEnforceSubscriptionsCron, paymentsproto.ActorPublishReminderCron:
		return true, nil
	}

	return false, nil
}

func formatCronitorMsg(job, actor, status string, timestamp time.Time) string {
	header := ":shark:    `CRONITOR: PHIL MITCHELL`    :robot:"
	base := `
Job: %s
Status: %s
Actor: %s
Timestamp: %v
	`
	formattedBase := fmt.Sprintf(base, job, status, actor, timestamp)
	return fmt.Sprintf("%s```%s```", header, formattedBase)
}
