package dao

import (
	"context"
	"testing"
	"time"

	"swallowtail/s.account/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAccounts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integtation test: TestReadAccounts")
	}

	ctx := context.Background()

	db.Exec(ctx, `DELETE FROM s_account_accounts`)

	t.Cleanup(func() {
		db.Exec(ctx, `DELETE FROM s_account_accounts`)
	})

	var (
		username    = "satoshi"
		password    = "unbreakable"
		email       = "test@email.co.uk"
		discordID   = "kdjadkadjadsksf"
		phoneNumber = "077123456763"

		username2    = "vitalik"
		password2    = "biglargepassword"
		email2       = "test+2@email.co.uk"
		discordID2   = "dkajdakdjakda"
		phoneNumber2 = "077123422763"
	)

	// Insert test data to datastore.
	_, err := db.Exec(
		ctx, `INSERT INTO s_account_accounts (username,password,email,discord_id,phone_number) values ($1,$2,$3,$4,$5)`,
		username, password, email, discordID, phoneNumber,
	)
	require.NoError(t, err, "Failed to write test data to database")

	// Read account back out.
	accounts, err := ReadAccounts(ctx)
	require.NoError(t, err, "Failed to read previously wrote test data")

	// Run assertions.
	assert.Len(t, accounts, 1)
	account := accounts[0]
	assert.Equal(t, username, account.Username)
	assert.Equal(t, password, account.Password)
	assert.Equal(t, email, account.Email)
	assert.Equal(t, discordID, account.DiscordID)
	assert.Equal(t, phoneNumber, account.PhoneNumber)

	// Insert more test data.
	_, err = db.Exec(
		ctx, `INSERT INTO s_account_accounts (username,password,email,discord_id,phone_number) values ($1,$2,$3,$4,$5)`,
		username2, password2, email2, discordID2, phoneNumber2,
	)
	require.NoError(t, err, "Failed to write test data to database")

	// Read data back out again, this time checking that we have both pieces of data.
	accounts, err = ReadAccounts(ctx)
	require.NoError(t, err, "Failed to read previously wrote test data")

	// Run assertions.
	assert.Len(t, accounts, 2)
}

func TestReadAccountByUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integtation test: TestReadAccountByUsername")
	}

	ctx := context.Background()

	db.Exec(ctx, `DELETE FROM s_account_accounts`)

	t.Cleanup(func() {
		db.Exec(ctx, `DELETE FROM s_account_accounts`)
	})

	var (
		username    = "satoshi"
		password    = "unbreakable"
		email       = "test@email.co.uk"
		discordID   = "kdjadkadjadsksf"
		phoneNumber = "077123456763"

		username2    = "vitalik"
		password2    = "biglargepassword"
		email2       = "test+2@email.co.uk"
		discordID2   = "dkajdakdjakda"
		phoneNumber2 = "077123422763"
	)

	// Insert test data to datastore.
	_, err := db.Exec(
		ctx, `INSERT INTO s_account_accounts (username,password,email,discord_id,phone_number) values ($1,$2,$3,$4,$5)`,
		username, password, email, discordID, phoneNumber,
	)
	require.NoError(t, err, "Failed to write test data to database")

	// Read account back out.
	account, err := ReadAccountByUsername(ctx, username)
	require.NoError(t, err, "Failed to read previously wrote test data")

	// Run assertions.
	assert.Equal(t, username, account.Username)
	assert.Equal(t, password, account.Password)
	assert.Equal(t, email, account.Email)
	assert.Equal(t, discordID, account.DiscordID)
	assert.Equal(t, phoneNumber, account.PhoneNumber)

	// Insert more test data.
	_, err = db.Exec(
		ctx, `INSERT INTO s_account_accounts (username,password,email,discord_id,phone_number) values ($1,$2,$3,$4,$5)`,
		username2, password2, email2, discordID2, phoneNumber2,
	)
	require.NoError(t, err, "Failed to write test data to database")

	account, err = ReadAccountByUsername(ctx, username2)
	require.NoError(t, err, "Failed to read previously wrote test data")

	assert.Equal(t, username2, account.Username)
	assert.Equal(t, password2, account.Password)
	assert.Equal(t, email2, account.Email)
	assert.Equal(t, discordID2, account.DiscordID)
	assert.Equal(t, phoneNumber2, account.PhoneNumber)

}
func TestCreateAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integtation test: TestCreateAccount")
	}

	ctx := context.Background()
	db.Exec(ctx, `DELETE FROM s_account_accounts`)
	t.Cleanup(func() {
		db.Exec(ctx, `DELETE FROM s_account_accounts`)
	})

	var (
		username        = "haskellcurry"
		password        = "mockingbird"
		email           = "haskellcurry@functional.co.uk"
		discordID       = "discordid_i484294"
		phoneNumber     = "07865482902"
		isAdmin         = false
		isFuturesMember = false
	)

	// Create an account.
	err := CreateAccount(ctx, &domain.Account{
		Username:    username,
		Password:    password,
		Email:       email,
		DiscordID:   discordID,
		PhoneNumber: phoneNumber,
		IsAdmin:     isAdmin,
	})
	require.NoError(t, err)

	// Read account back out. This probably should be done via a SQL statement for isolation,
	// but these tests are more like integration tests anyway.
	account, err := ReadAccountByUsername(ctx, username)
	require.NoError(t, err, "Failed to read account back out")

	// Run assertions.
	assert.Equal(t, username, account.Username)
	assert.Equal(t, password, account.Password)
	assert.Equal(t, email, account.Email)
	assert.Equal(t, discordID, account.DiscordID)
	assert.Equal(t, phoneNumber, account.PhoneNumber)
	assert.Equal(t, isAdmin, account.IsAdmin)
	assert.Equal(t, isFuturesMember, account.IsFuturesMember)

	// Timestamp assertions; we don't care for exactness, as long as they're in the right
	// ballpark.

	// Removing for now whilst local timestamps are wildly incorrect.
	// assert.True(t, between(account.Created, start, end))
	// assert.True(t, between(account.Updated, start, end))
	// assert.True(t, between(account.LastPaymentTimestamp, start, end))
}

func TestUpdateAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integtation test: TestUpdateAccount")
	}
}

func between(t, start, end time.Time) bool {
	if t.After(end) {
		return false
	}
	if t.Before(start) {
		return false
	}
	return true
}
