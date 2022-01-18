package handler

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.payments/dao"
	paymentsproto "swallowtail/s.payments/proto"
)

// ReadUsersLastPayment ...
func (s *PaymentsService) ReadUsersLastPayment(
	ctx context.Context, in *paymentsproto.ReadUsersLastPaymentRequest,
) (*paymentsproto.ReadUsersLastPaymentResponse, error) {
	// Validation.
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"actor_id": in.ActorId,
	}

	// Validation on the actor & request context.
	if err := isActorValid(in.ActorId); err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_users_last_payment_timestamp.unauthorized", errParams)
	}

	// Read from persistance layer.
	ts, err := dao.ReadUsersLastPaymentTimestamp(ctx, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_users_last_payment_timestamp", errParams)
	}

	return &paymentsproto.ReadUsersLastPaymentResponse{
		LastPaymentTimestamp:    timestamppb.New(*ts),
		HasUserPaidForLastMonth: currentMonthStartFromTimestamp(*ts) == currentMonthStartFromTimestamp(time.Now().UTC()),
	}, nil
}
