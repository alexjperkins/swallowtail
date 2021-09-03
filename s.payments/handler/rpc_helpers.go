package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"
	ftxproto "swallowtail/s.ftx/proto"
)

// setUserAsFuturesMember sets the user as a futures member.
func setUserAsFuturesMember(ctx context.Context, userID string) error {
	// Set user as a futures member internally.
	_, err := (&accountproto.UpdateAccountRequest{
		ActorId:   accountproto.ActorSystemPayments,
		UserId:    userID,
		IsFutures: true,
	}).Send(ctx).Response()
	if err != nil {
		return gerrors.Augment(err, "failed_to_set_user_as_futures_member", map[string]string{
			"user_id": userID,
		})
	}

	// Set user as discord futures member.
	if _, err := (&discordproto.UpdateUserRolesRequest{
		ActorId:           discordproto.DiscordRolesUpdateActorPaymentsSystem,
		UserId:            userID,
		MergeWithExisting: true,
		Roles: []*discordproto.Role{
			{
				RoleId:   discordproto.DiscordSatoshiFuturesRoleID,
				RoleName: discordproto.DiscordSatoshiFuturesRole,
			},
		},
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_set_user_as_futures_member", nil)
	}

	return nil
}

// readUserRoles reads the users roles from discord.
func readUserRoles(ctx context.Context, userID string) ([]*discordproto.Role, error) {
	rsp, err := (&discordproto.ReadUserRolesRequest{
		UserId: userID,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_user_discord_roles", map[string]string{
			"user_id": userID,
		})
	}

	return rsp.GetRoles(), nil
}

// isUserRegistered checks if the userID has an account registered in `s.account`
func isUserRegistered(ctx context.Context, userID string) (*accountproto.Account, error) {
	rsp, err := (&accountproto.ReadAccountRequest{
		UserId: userID,
	}).Send(ctx).Response()
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		return nil, nil
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_check_if_user_register", nil)
	}

	return rsp.Account, nil
}

// isMonthlyTransactionInDepositAccount checks if the txid exists in the deposit account.
func isMonthlyTransactionInDepositAccount(ctx context.Context, transactionID string, minimumExpectedAmount float64) (bool, error) {
	start := currentMonthStartTimestamp()

	rsp, err := (&ftxproto.ListAccountDepositsRequest{
		ActorId: ftxproto.FTXDepositAccountActorPaymentsSystem,
		// We require only second granularity.
		Start: start.Unix(),
	}).Send(ctx).Response()
	if err != nil {
		return false, gerrors.Augment(err, "failed_to_check_txid_in_deposit_account", map[string]string{
			"txid": transactionID,
		})
	}

	thisMonthsDeposits := rsp.Deposits

	for _, deposit := range thisMonthsDeposits {
		if deposit.TransactionId == transactionID {
			return float64(deposit.Size) > minimumExpectedAmount, nil
		}
	}

	return false, nil
}
