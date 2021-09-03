package dao

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.payments/domain"
	"time"
)

func ReadPaymentByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	var (
		sql = `
		SELECT * FROM s_payments_payments
		WHERE transaction_id=$1
		`
		payments []*domain.Payment
	)

	if err := db.Select(ctx, &payments, sql); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(payments) {
	case 0:
		return nil, nil
	case 1:
		return payments[0], nil
	default:
		return nil, gerrors.FailedPrecondition("read_payment_by_transaction_id.multiple_unique_payments", nil)
	}
}

// RegisterPayment ...
func RegisterPayment(ctx context.Context, payment *domain.Payment) error {
	var (
		sql = ` 
		INSERT INTO s_payments_payments(
			user_id, transaction_id, timestamp, amount_in_usdt, audit_note
		)
		VALUES
			$1, $2, $3, $4, $5
		`
	)

	if _, err := (db.Exec(
		ctx, sql,
		payment.UserID, payment.TransactionID, payment.Timestamp, payment.AmountInUSDT, payment.AuditNote,
	)); err != nil {
		return gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil
}

// UserPaymentExistsSince ...
func UserPaymentExistsSince(ctx context.Context, userID string, after time.Time) (bool, error) {
	var (
		sql = `
		SELECT * FROM s_payments_payments
		WHERE user_id=$1
		AND
		timestamp >= $2
		`
		payments []*domain.Payment
	)

	if err := db.Select(ctx, &payments, sql); err != nil {
		return false, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(payments) {
	case 0:
		return false, nil
	default:
		return true, nil
	}
}
