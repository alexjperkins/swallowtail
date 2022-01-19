package handler

import (
	"context"
	"strconv"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.payments/dao"
	"swallowtail/s.payments/marshaling"
	paymentsproto "swallowtail/s.payments/proto"
)

// ListPaymentsByUserID ...
func (s *PaymentsService) ListPaymentsByUserID(
	ctx context.Context, in *paymentsproto.ListPaymentsByUserIDRequest,
) (*paymentsproto.ListPaymentsByUserIDResponse, error) {
	// Validate.
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	case in.Limit == 0:
		return nil, gerrors.BadParam("missing_param.limit", nil)
	case in.Limit > 120:
		return nil, gerrors.FailedPrecondition("limit_too_high", nil)
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"actor_id": in.ActorId,
		"limit":    strconv.Itoa(int(in.Limit)),
	}

	// Authenticate actor.
	if err := isActorValid(in.ActorId); err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_payments_by_user_id", errParams)
	}

	// Read from persistance layer.
	payments, err := dao.ListPaymentsByUserID(ctx, in.UserId, int(in.Limit))
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_payments_by_user_id", errParams)
	}

	// Marshal to proto.
	protos := marshaling.PaymentsToProtos(payments)

	return &paymentsproto.ListPaymentsByUserIDResponse{
		Payments: protos,
	}, nil
}
