package handler

import (
	"context"
	"fmt"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"
	ftxproto "swallowtail/s.ftx/proto"
	"time"
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

// removeUserAsFuturesMember sets the user as a futures member.
func removeUserAsFuturesMember(ctx context.Context, userID string) error {
	// Remove user as a futures member internally.
	_, err := (&accountproto.UpdateAccountRequest{
		ActorId:   accountproto.ActorSystemPayments,
		UserId:    userID,
		IsFutures: false,
	}).Send(ctx).Response()
	if err != nil {
		return gerrors.Augment(err, "failed_to_remove_user_as_futures_member", nil)
	}

	if _, err := (&discordproto.RemoveUserRoleRequest{
		ActorId: discordproto.DiscordRolesUpdateActorPaymentsSystem,
		UserId:  userID,
		Role: &discordproto.Role{
			RoleId:   discordproto.DiscordSatoshiFuturesRoleID,
			RoleName: discordproto.DiscordSatoshiFuturesRole,
		},
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_remove_user_as_discord_futures_member", nil)
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

func listFuturesMembers(ctx context.Context) ([]*accountproto.Account, error) {
	rsp, err := (&accountproto.ListAccountsRequest{
		IsFuturesMember: true,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_futures_members", nil)
	}

	return rsp.Accounts, nil
}

func postToPaymentsPulseChannel(ctx context.Context, isExistingFuturesMember bool, userID, username, transactionID, auditNote string, amount float64, timestamp time.Time) error {
	header := ":money_with_wings:   `PAYMENT RECEIVED`   :money_with_wings:"
	content := `
UserID: %s
Username: %s
TXID: %s
AuditNote: %s
AmountInUSDT: %d
IsExistingMember: %v
Timestamp: %v
	`
	formattedContent := fmt.Sprintf(content, userID, username, transactionID, auditNote, amount, isExistingFuturesMember, timestamp)

	// Best Effort
	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiPaymentsPulseChannel,
		SenderId:       "system-payments",
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: formattedContent,
	}).Send(ctx).Response()
	if err != nil {
		return gerrors.Augment(err, "failed_to_post_to_payments_pulse_channel", nil)
	}

	return nil
}

func postToAccountsPulseChannel(ctx context.Context, isExistingFuturesMember bool, userID, username string, timestamp time.Time) error {
	if isExistingFuturesMember {
		// If the user is already a futures member then we don't have a new member; in such case theres no point posting to
		// out pulse channel.
		return nil
	}

	header := ":money_mouth:    New Futures Member    :face_with_monocle"
	content := `
UserID: %s
Username: %s
Timestamp: %v
	`
	formattedContent := fmt.Sprintf(content, userID, username, timestamp)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiAccountsPulseChannel,
		SenderId:       "system-payments",
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: formattedContent,
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_post_to_account_pulse_channel", nil)
	}

	return nil
}
