package dao

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/domain"
	"time"

	"github.com/imdario/mergo"
	"github.com/monzo/slog"
)

// ReadVenueAccountByVenueAccountID ...
func ReadVenueAccountByVenueAccountID(ctx context.Context, venueAccountID string) (*domain.VenueAccount, error) {
	var (
		sql = `
		SELECT * FROM s_account_venue_accounts
		WHERE venue_account_id=$1
		`
		venueAccounts []*domain.VenueAccount
	)

	if err := db.Select(ctx, venueAccounts, sql, venueAccountID); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	if len(venueAccounts) == 0 {
		return nil, gerrors.NotFound("venue_account_not_found", nil)
	}

	return venueAccounts[0], nil
}

// ReadVenueAccountByVenueAccountDetails ...
func ReadVenueAccountByVenueAccountDetails(ctx context.Context, venueID, userID, subaccount string) (*domain.VenueAccount, error) {
	var (
		baseSql = `
		SELECT * FROM s_account_venue_accounts
		WHERE venue_id=$1
		AND user=$2
		`
		venueAccounts []*domain.VenueAccount
	)

	var sql = baseSql
	if subaccount != "" {
		sql = baseSql + `AND subaccount=$3`
	}

	if err := db.Select(ctx, &venueAccounts, sql, strings.ToUpper(venueID), userID, subaccount); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(venueAccounts) {
	case 0:
		return nil, gerrors.NotFound("venue_accounts_not_found_for_user_id", nil)
	case 1:
		return venueAccounts[0], nil
	default:
		slog.Critical(ctx, "Inconsistent state: more than one identical venue account found for user", map[string]string{
			"venue_id":   venueID,
			"user_id":    userID,
			"subaccount": subaccount,
		})
		return venueAccounts[0], nil
	}
}

// ListVenueAccountsByUserID ...
func ListVenueAccountsByUserID(ctx context.Context, userID string, isActive bool) ([]*domain.VenueAccount, error) {
	var (
		sql = `
		SELECT * FROM s_account_venue_accounts
		WHERE user_id=$1
		AND is_active=$2
		ORDER BY venue_id 
		`
		venueAccounts []*domain.VenueAccount
	)

	if err := db.Select(ctx, &venueAccounts, sql, userID, isActive); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	if len(venueAccounts) == 0 {
		return nil, gerrors.NotFound("venue_accounts_not_found_for_user_id", nil)
	}

	return venueAccounts, nil
}

// AddVenueAccount ...
func AddVenueAccount(ctx context.Context, venueAccount *domain.VenueAccount) error {
	var (
		sql = `
		INSERT INTO s_account_venue_accounts
			(venue_id, api_key, secret_key, user_id, created, updated, is_active, account_alias)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, &8)
		`
	)

	now := time.Now().UTC()
	if _, err := (db.Exec(
		ctx, sql,
		venueAccount.VenueID, venueAccount.APIKey, venueAccount.SecretKey, venueAccount.UserID,
		now, now,
		venueAccount.IsActive,
		venueAccount.AccountAlias,
	)); err != nil {
		return gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil
}

// RemoveVenueAccount ...
func RemoveVenueAccount(ctx context.Context, venueAccountID string) error {
	var (
		sql = `
		DELETE FROM s_account_venue_accounts
		WHERE venue_account_id=$1 
		`
	)

	if _, err := db.Exec(
		ctx, sql, venueAccountID,
	); err != nil {
		return gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}
	return nil
}

// UpdateVenueAccount ...
func UpdateVenueAccount(ctx context.Context, mutation *domain.VenueAccount) (*domain.VenueAccount, error) {
	var (
		sql = `
		UPDATE s_account_venue_accounts
		SET
			api_key=$2, secret_key=$3, account_alias=$4 updated=$5
		WHERE venue_account_id=$1
		`
	)
	if mutation.VenueAccountID == "" {
		return nil, gerrors.FailedPrecondition("missing_venue_account_id", nil)
	}

	venueAccount, err := ReadVenueAccountByVenueAccountID(ctx, mutation.VenueAccountID)
	if err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	if err := mergo.Merge(&venueAccount, mutation); err != nil {
		return nil, gerrors.Augment(err, "failed_to_merge_venue_account_update_request", nil)
	}

	venueAccount.Updated = time.Now().UTC()

	if _, err := db.Exec(
		ctx, sql,
		venueAccount.VenueAccountID, venueAccount.APIKey, venueAccount.SecretKey, venueAccount.AccountAlias, venueAccount.Updated,
	); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil, nil
}
