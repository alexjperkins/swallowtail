package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.payments/dao"
	paymentsproto "swallowtail/s.payments/proto"

	"github.com/monzo/slog"
)

// EnforceSubscriptions ...
func (s *PaymentsService) EnforceSubscriptions(
	ctx context.Context, in *paymentsproto.EnforceSubscriptionRequest,
) (*paymentsproto.EnforceSubscriptionResponse, error) {
	// TODO validate the day.
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

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

		ok, err := dao.UserPaymentExistsSince(ctx, fm.UserId, currentMonthStartTimestamp())
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions", errParams)
		}

		// Futures member has paid for the month; we can continue.
		if ok {
			continue
		}

		slog.Info(ctx, "Offboarding user: %v: %v", fm.Username, fm.UserId)

		// Uh-oh they haven't paid, lets offboard them.
		if err := offboardSubscriber(ctx, fm.UserId); err != nil {
			return nil, gerrors.Augment(err, "failed_to_enforce_subscriptions", errParams)
		}
	}

	return &paymentsproto.EnforceSubscriptionResponse{}, nil
}

func isActorValid(ctx context.Context, actorID string) (bool, error) {
	return false, nil
}
