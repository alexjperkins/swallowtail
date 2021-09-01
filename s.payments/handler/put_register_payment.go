package handler

import (
	"context"
	"strconv"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.payments/dao"
	paymentsproto "swallowtail/s.payments/proto"
)

// RegisterPayment ...
func (s *PaymentsService) RegisterPayment(
	ctx context.Context, in *paymentsproto.RegisterPaymentRequest,
) (*paymentsproto.RegisterPaymentResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.TransactionId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.AmountInUsdt == 0:
		return nil, gerrors.BadParam("bad_param.amount_in_usdt_cannot_be_zero", nil)
	case in.AmountInUsdt < 0:
		return nil, gerrors.BadParam("bad_param.amount_in_usdt_cannot_be_negative", nil)
	}

	errParams := map[string]string{
		"user_id":        in.UserId,
		"transaction_id": in.TransactionId,
		"amount_in_usdt": strconv.FormatFloat(float64(in.AmountInUsdt), 'f', 6, 64),
	}

	// Check the user does indeed have an account.
	ok, err := isUserRegistered(ctx, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_register_payment", errParams)
	}

	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_register_payment.user_does_not_have_an_account", errParams)
	}

	// Check that the txid doesn't already exist.
	payment, err := dao.ReadPaymentByTransactionID(ctx, in.TransactionId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_register_payment", errParams)
	}

	if payment != nil {
		return nil, gerrors.AlreadyExists("failed_to_register_payment.payment_already_exists", errParams)
	}

	// Check the user hasn't already paid this month
	hasAlreadyPaid, err := dao.UserPaymentExistsSince(ctx, in.UserId, currentMonthStartTimestamp())
	if err != nil {
		return nil, gerrors.AlreadyExists("failed_to_register_payment.failed_check_if_user_already_paid", errParams)
	}

	if hasAlreadyPaid {
		return nil, gerrors.FailedPrecondition("failed_to_register_payment.user_has_already_paid", errParams)
	}

	// TODO: validate txid regex.
	// TODO: validate txid & the amount on FTX.
	return nil, nil
}
