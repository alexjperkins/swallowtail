package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.payments/dao"
	"swallowtail/s.payments/domain"
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

	now := time.Now().UTC()
	timestampOfChecks := currentMonthStartFromTimestamp(now)

	errParams := map[string]string{
		"user_id":             in.UserId,
		"transaction_id":      in.TransactionId,
		"amount_in_usdt":      strconv.FormatFloat(float64(in.AmountInUsdt), 'f', 6, 64),
		"timestamp_of_checks": timestampOfChecks.String(),
	}

	// Check the user does indeed have an account.
	account, err := readAccount(ctx, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_register_payment", errParams)
	}

	// Check that the txid doesn't already exist.
	payment, err := dao.ReadPaymentByTransactionID(ctx, in.TransactionId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "payment_not_found"):
	case payment != nil:
		return nil, gerrors.Augment(err, "payment_already_made", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_register_payment.read_payment", errParams)
	}

	// Check the user hasn't already paid this month
	hasAlreadyPaid, err := dao.UserPaymentExistsSince(ctx, in.UserId, timestampOfChecks)
	if err != nil {
		slog.Info(ctx, "User has already paid: timestamp: %v user_id: %s err: %v", timestampOfChecks, in.UserId, err)
		return nil, gerrors.Augment(err, "failed_to_register_payment.failed_check_if_user_already_paid", errParams)
	}
	if hasAlreadyPaid {
		return nil, gerrors.FailedPrecondition("failed_to_register_payment.user_has_already_paid", errParams)
	}

	// We check if the monthly amount has indeed been paid; we allow a discrepancy of 1 to allow for tx fees.
	doesTxExist, err := isMonthlyTransactionInDepositAccount(ctx, in.TransactionId, float64(in.AmountInUsdt)-1)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_register_payment", errParams)
	}
	if !doesTxExist {
		return nil, gerrors.FailedPrecondition("failed_to_register_payment.transaction_of_correct_amount_does_not_exist_in_deposit_account", errParams)
	}

	// Set user as a futures member on s.account & in discord
	if err := setUserAsFuturesMember(ctx, in.UserId); err != nil {
		return nil, gerrors.Augment(err, "failed_to_register_payment.failed_to_set_user_as_futures_member", errParams)
	}

	slog.Info(ctx, "User: %s, set as a futures member", in.UserId)

	// Okay; everything is in check, we can now safely store the tx to our persistence layer.
	//
	// NOTE: We push to our persistence layer since all the above are idempotent. If we fail to store before
	// setting the user as a futures maybe, then on retry we fail since the txid will already exist.
	if err := dao.RegisterPayment(ctx, &domain.Payment{
		UserID:        in.UserId,
		TransactionID: in.TransactionId,
		AuditNote:     in.AuditNote,
		AmountInUSDT:  float64(in.AmountInUsdt),
		Timestamp:     now,
	}); err != nil {
		return nil, gerrors.Augment(err, "failed_to_register_payment.persistence_layer", errParams)
	}

	slog.Info(ctx, "Payment registered for user: %s", in.UserId)

	// Best effort; post to pulse channels
	if err := postToPaymentsPulseChannel(ctx, account.IsFuturesMember, in.UserId, account.Username, in.TransactionId, in.AuditNote, float64(in.AmountInUsdt), now); err != nil {
		slog.Error(ctx, "Failed to publish to payments pulse channel: %v: Error", in.UserId, err)
	}

	if err := postToAccountsPulseChannel(ctx, account.IsFuturesMember, in.UserId, account.Username, now); err != nil {
		slog.Error(ctx, "Failed to publish to accounts pulse channel: %v: Error", in.UserId, err)
	}

	return &paymentsproto.RegisterPaymentResponse{}, nil
}
