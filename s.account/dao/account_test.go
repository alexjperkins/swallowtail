package dao

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAccounts(t *testing.T) {
	ctx := context.Background()

	db.Exec(ctx, `DELETE FROM accounts`)

	t.Cleanup(func() {
		db.Exec(ctx, `DELETE FROM accounts`)
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
		ctx, `INSERT INTO accounts (username,password,email,discord_id,phone_number) values ($1,$2,$3,$4,$5)`,
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
		ctx, `INSERT INTO accounts (username,password,email,discord_id,phone_number) values ($1,$2,$3,$4,$5)`,
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
	ctx := context.Background()

	db.Exec(ctx, `DELETE FROM accounts`)

	t.Cleanup(func() {
		db.Exec(ctx, `DELETE FROM accounts`)
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
		ctx, `INSERT INTO accounts (username,password,email,discord_id,phone_number) values ($1,$2,$3,$4,$5)`,
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
		ctx, `INSERT INTO accounts (username,password,email,discord_id,phone_number) values ($1,$2,$3,$4,$5)`,
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
func TestCreateAccount(t *testing.T) {}
func TestUpdateAccount(t *testing.T) {}
