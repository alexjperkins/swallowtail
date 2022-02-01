package handler

import (
	"context"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.payments/dao"
	paymentsproto "swallowtail/s.payments/proto"
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
	if err := isActorValid(in.ActorId); err != nil {
		return nil, gerrors.Unauthenticated("failed_to_enforce_subscriptions", errParams)
	}

	// Notify Pulse channel.
	if _, err := (&discordproto.SendMsgToChannelRequest{
		Content:   formatCronitorMsg("Enforce Subscriptions", in.ActorId, "Started", time.Now().Truncate(time.Second)),
		ChannelId: discordproto.DiscordSatoshiPaymentsPulseChannel,
		SenderId:  "system:payments",
		Force:     true,
	}).Send(ctx).Response(); err != nil {
		slog.Error(ctx, "Failed to notify pulse channel of subsciption enforcement")
	}

	futuresMembers, err := listFuturesMembers(ctx)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions", errParams)
	}

	now := time.Now().UTC()
	errParams["enforcement_timestamp"] = now.String()

	slog.Info(ctx, "Subscription enforcement starting: total futures members: %d", len(futuresMembers))

	for _, fm := range futuresMembers {
		if fm.IsAdmin {
			slog.Warn(ctx, "Skipping subscription payment check for admin: %s: %s", fm.UserId, fm.Username)
			continue
		}

		errParams["user_id"] = fm.UserId

		ok, err := dao.UserPaymentExistsSince(ctx, fm.UserId, currentMonthStartFromTimestamp(now))
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions", errParams)
		}
		if ok {
			// Futures member has paid for the month; we can continue.
			continue
		}

		slog.Info(ctx, "Offboarding user: %v: %v", fm.Username, fm.UserId)

		// Uh-oh they haven't paid, lets offboard them.
		if err := offboardSubscriber(ctx, fm.UserId, fm.Username); err != nil {
			slog.Error(ctx, "Failed to enforce subscriber, Error: %v", err)
			return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions", errParams)
		}
	}

	return &paymentsproto.EnforceSubscriptionsResponse{}, nil
}
