package dao

import (
	"context"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/sql"
	"swallowtail/s.payments/domain"
)

// ReadPaymentByTransactionID ...
func ReadPaymentByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	var (
		sql = `
		SELECT * FROM s_payments_payments
		WHERE transaction_id=$1
		`
		payments []*domain.Payment
	)

	if err := db.Select(ctx, &payments, sql, transactionID); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(payments) {
	case 0:
		return nil, gerrors.NotFound("payment_not_found", nil)
	case 1:
		return payments[0], nil
	default:
		slog.Critical(ctx, "Incoherent state of persistance layer: multiple txid of unique payment records: %v", transactionID)
		return nil, gerrors.FailedPrecondition("read_payment_by_transaction_id.multiple_unique_payments", nil)
	}
}

// RegisterPayment ...
func RegisterPayment(ctx context.Context, payment *domain.Payment) error {
	var (
		sql = ` 
		INSERT INTO s_payments_payments(
			user_id, transaction_id, payment_timestamp, amount_in_usdt, audit_note
		)
		VALUES (
			$1, $2, $3, $4, $5
		)`
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
		SELECT EXISTS (
			SELECT 1 FROM s_payments_payments
			WHERE user_id=$1
			AND payment_timestamp >= $2
		)`
		hasPaid bool
	)

	if err := db.Get(ctx, &hasPaid, sql, userID, after); err != nil {
		return false, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return hasPaid, nil
}

// ReadUsersLastPaymentTimestamp ...
func ReadUsersLastPaymentTimestamp(ctx context.Context, userID string) (*time.Time, error) {
	var (
		query = `
		SELECT payment_timestamp FROM s_payments_payments
		WHERE user_id=$1
		ORDER BY payment_timestamp DESC
		LIMIT 1
		`
		lastPaymentTimestamp time.Time
	)

	if err := db.Get(ctx, &lastPaymentTimestamp, query, userID); err != nil {
		return nil, sql.PostgresGetFailed(err)
	}

	return &lastPaymentTimestamp, nil
}

// ListPaymentsByUserID ...
func ListPaymentsByUserID(ctx context.Context, userID string, limit int) ([]*domain.Payment, error) {
	var (
		query = `
		SELECT 
			transaction_id,
			payment_timestamp,
			amount_in_usdt
		FROM s_payments_payments
		WHERE user_id=$1
		ORDER BY payment_timestamp DESC
		LIMIT $2
		`
		payments []*domain.Payment
	)

	if err := db.Select(ctx, &payments, query, userID, limit); err != nil {
		return nil, sql.PostgresSelectFailed(err)
	}

	return payments, nil
}
