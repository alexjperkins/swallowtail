package dao

import (
	"context"
	"swallowtail/s.account/domain"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/imdario/mergo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

// ReadAccounts returns a list of all domain accounts from the underlying datastore.
func ReadAccounts(ctx context.Context) ([]*domain.Account, error) {
	var (
		sql      = `SELECT * FROM accounts`
		accounts []*domain.Account
	)
	err := pgxscan.Select(ctx, db, &accounts, sql)
	if err != nil {
		return nil, terrors.Propagate(err)
	}
	return accounts, nil
}

// ReadAccountByID returns the a domain account from the underlying datastore, by ID.
func ReadAccountByID(ctx context.Context, accountID string) (*domain.Account, error) {
	var (
		sql      = `SELECT * FROM accounts WHERE account_id=$1`
		accounts []*domain.Account
	)
	err := pgxscan.Select(ctx, db, &accounts, sql, accountID)
	if err != nil {
		return nil, terrors.Propagate(err)
	}
	if len(accounts) > 1 {
		// We shouldn't have this case, since this is enforced at the persistence layer.
		slog.Error(ctx, "More than one account with identical id", map[string]string{
			"account_id": accountID,
		})
	}
	return accounts[0], nil
}

// ReadAccountByUsername returns the a domain account from the underlying datastore, by username.
func ReadAccountByUsername(ctx context.Context, username string) (*domain.Account, error) {
	var (
		sql      = `SELECT * FROM accounts WHERE username=$1`
		accounts []*domain.Account
	)
	err := pgxscan.Select(ctx, db, &accounts, sql, username)
	if err != nil {
		return nil, terrors.Propagate(err)
	}
	if len(accounts) > 1 {
		// We shouldn't have this case, since this is enforced at the persistence layer.
		slog.Error(ctx, "More than one account with identical username", map[string]string{
			"username": username,
		})
	}
	return accounts[0], nil
}

// CreateAccount creates a new account in the datastore.
func CreateAccount(ctx context.Context, account *domain.Account) error {
	var (
		sql = `INSERT INTO accounts(username,password email,discord_id,phone_number,high_priority_pager,low_priority_pager) values($1,$2,$3,$4,$5,$6,$7)`
	)
	if _, err := (db.Exec(
		ctx, sql,
		account.Email, account.Password, account.Email, account.DiscordID, account.PhoneNumber, account.HighPriorityPager,
		account.LowPriorityPager,
	)); err != nil {
		return terrors.Propagate(err)
	}
	return nil
}

// UpdateAccount updates an already existing account; it will perform a mutation on the existing account.
// AccountID must be provided at least to the passed domain account struct, `mutation`.
func UpdateAccount(ctx context.Context, mutation *domain.Account) (*domain.Account, error) {
	var (
		sql = `UPDATE accounts set username=$1,password=$2,email=$3,discord_id=$4,phone_number$5,high_priority_pager=$6,low_priority_pager=$7 WHERE user_id=$8`
	)
	if mutation.AccountID == "" {
		return nil, terrors.PreconditionFailed("mutation-without-id", "Account mutation requires at least the account ID", nil)
	}

	account, err := ReadAccountByID(ctx, mutation.AccountID)
	if err != nil {
		return nil, terrors.Propagate(err)
	}

	if err := mergo.Merge(&account, mutation); err != nil {
		return nil, terrors.BadRequest("mutation-merge-failure", "Failed to merge account mutation", map[string]string{
			"upstream_err": err.Error(),
		})
	}

	if _, err := (db.Exec(
		ctx, sql,
		account.Email, account.Password, account.Email, account.DiscordID, account.PhoneNumber, account.HighPriorityPager,
		account.LowPriorityPager, account.AccountID,
	)); err != nil {
		return nil, terrors.Propagate(err)
	}
	return account, nil
}
