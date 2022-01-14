package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/domain"
)

// ReadVenueAccountByAccountAlias ...
func ReadVenueAccountByAccountAlias(ctx context.Context, userID, accountAlias string) (*domain.VenueAccount, error) {
	var (
		sql = `
		SELECT 
			venue_account_id,
			venue_id,
			api_key,
			secret_key,
			subaccount,
			user_id,
			created,
			updated,
			is_active,
			account_alias,
			COALESCE(url, '') as url,
			COALESCE(ws_url, '') as ws_url
		FROM s_account_venue_accounts
		WHERE
			user_id=$1
		AND
			account_alias=$2
		`
		venueAccounts []*domain.VenueAccount
	)

	if err := db.Select(ctx, venueAccounts, sql, userID, accountAlias); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(venueAccounts) {
	case 0:
		return nil, gerrors.NotFound("venue_account_not_found", nil)
	case 1:
		return venueAccounts[0], nil
	default:
		slog.Critical(ctx, "Incoherent persistance state: violation of unique constraint: (user_id, account_alias)")
		return venueAccounts[0], nil
	}
}

// ReadVenueAccountByVenueAccountID ...
func ReadVenueAccountByVenueAccountID(ctx context.Context, venueAccountID string) (*domain.VenueAccount, error) {
	var (
		sql = `
		SELECT
			venue_account_id,
			venue_id,
			api_key,
			secret_key,
			subaccount,
			user_id,
			created,
			updated,
			is_active,
			account_alias,
			COALESCE(url, '') as url,
			COALESCE(ws_url, '') as ws_url
		FROM s_account_venue_accounts
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
func ReadVenueAccountByVenueAccountDetails(ctx context.Context, venueID, userID, _ string) (*domain.VenueAccount, error) {
	var (
		sql = `
		SELECT 
			venue_account_id,
			venue_id,
			api_key,
			secret_key,
			subaccount,
			user_id,
			created,
			updated,
			is_active,
			account_alias,
			COALESCE(url, '') as url,
			COALESCE(ws_url, '') as ws_url
		FROM s_account_venue_accounts
		WHERE venue_id=$1
		AND user_id=$2
		AND account_alias=$3
		`
		venueAccounts []*domain.VenueAccount
	)

	// Switch empty subaccount to the default value we store in the db.
	//if subaccount == "" {
	//	subaccount = accountproto.SubAccountUnknown
	//}

	accountAlias := strings.ToUpper(fmt.Sprintf("%s-MAIN", venueID))

	if err := db.Select(ctx, &venueAccounts, sql, strings.ToUpper(venueID), userID, accountAlias); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(venueAccounts) {
	case 0:
		return nil, gerrors.NotFound("venue_account_not_found_for_user_id", nil)
	case 1:
		return venueAccounts[0], nil
	default:
		slog.Critical(ctx, "Inconsistent state: more than one identical venue account found for user", map[string]string{
			"venue_id": venueID,
			"user_id":  userID,
		})
		return venueAccounts[0], nil
	}
}

func ReadInternalVenueAccount(ctx context.Context, venueID, subaccount, internalAccountType string) (*domain.InternalVenueAccount, error) {
	var (
		sql = `
		SELECT 
			venue_account_id,
			venue_id,
			api_key,
			secret_key,
			subaccount,
			COALESCE(url, '') as url,
			COALESCE(ws_url, '') as ws_url,
			venue_account_type,
			created,
			updated
		FROM s_account_internal_venue_accounts
		WHERE 
			venue_id=$1
		AND
			subaccount=$2
		AND 
			venue_account_type=$3
		`
		internalAccounts []*domain.InternalVenueAccount
	)

	if err := db.Select(ctx, &internalAccounts, sql, venueID, subaccount, internalAccountType); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(internalAccounts) {
	case 0:
		return nil, gerrors.NotFound("internal_venue_account.not_found", nil)
	case 1:
		return internalAccounts[0], nil
	default:
		slog.Critical(ctx, "Unique constraint failed on s_account_internal_venue_accounts")
		return internalAccounts[0], nil
	}
}

// ListVenueAccountsByUserID ...
func ListVenueAccountsByUserID(ctx context.Context, userID string, isActive bool) ([]*domain.VenueAccount, error) {
	var (
		sql = `
		SELECT
			venue_account_id,
			venue_id,
			api_key,
			secret_key,
			subaccount,
			user_id,
			created,
			updated,
			is_active,
			account_alias,
			COALESCE(url, '') as url,
			COALESCE(ws_url, '') as ws_url
		FROM s_account_venue_accounts
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
			(venue_id, user_id, api_key, secret_key, subaccount, url, ws_url, account_alias, created, updated, is_active)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
	)

	// If we have an empty account alias then we build one internally using the venue & we
	// assume this is the main.
	// The unqiue constraint (user_id, account_alias) therefore forces the end user to
	// choose a user space unique account alias if they wish to have more than one venue account per venue.
	var accountAlias string
	switch {
	case venueAccount.AccountAlias == "":
		accountAlias = strings.ToUpper(fmt.Sprintf("%s-%s", venueAccount.VenueID, "MAIN"))
	}

	now := time.Now().UTC()
	if _, err := (db.Exec(
		ctx, sql,
		venueAccount.VenueID, venueAccount.UserID, venueAccount.APIKey, venueAccount.SecretKey, venueAccount.SubAccount,
		venueAccount.URL, venueAccount.WSURL, accountAlias,
		now, now,
		venueAccount.IsActive,
	)); err != nil {
		return gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil
}

// CreateOrUpdateInternalVenueAccount ...
func CreateOrUpdateInternalVenueAccount(ctx context.Context, venueAccount *domain.InternalVenueAccount, allowUpdate bool) error {
	var (
		sql = `
		INSERT INTO s_account_internal_venue_accounts
			(venue_id, api_key, secret_key, subaccount, url, ws_url, account_type, updated)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		`
	)

	// Convert insert into an upsert if we allow an update.
	if allowUpdate {
		sql = sql + `
		ON CONFLICT DO UPDATE
		`
	}

	if _, err := (db.Exec(
		ctx, sql,
		venueAccount.VenueID, venueAccount.APIKey, venueAccount.SecretKey, venueAccount.URL, venueAccount.WSURL, venueAccount.VenueAccountType,
		time.Now().UTC(),
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
