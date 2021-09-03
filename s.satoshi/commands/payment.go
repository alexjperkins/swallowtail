package commands

import (
	"context"
	"fmt"
	"swallowtail/libraries/gerrors"
	paymentsproto "swallowtail/s.payments/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	paymentCommandID = "payment"
	paymentUsage     = `
	Usage: !payment <subcommand>
	Description: command for automating payments & more.

	SubCommands:
	1. register: registers a new payment.
	`
)

func init() {
	register(paymentCommandID, &Command{
		ID:                  paymentCommandID,
		IsPrivate:           true,
		MinimumNumberOfArgs: 1,
		Usage:               paymentUsage,
		Handler:             paymentHandler,
		SubCommands: map[string]*Command{
			"register": {
				ID:                  "payment-register",
				IsPrivate:           true,
				MinimumNumberOfArgs: 1,
				Usage:               `!payment <transaction_id>`,
				Handler:             registerPaymentHandler,
				FailureMsg:          "Please check that you have an account registered; or maybe you've already paid for this month? ping @ajperkins if unsure",
			},
		},
	})
}

func paymentHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return gerrors.Unimplemented("parent_command_unimplemented.payment", nil)
}

func registerPaymentHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	txid := tokens[0]

	_, err := (&paymentsproto.RegisterPaymentRequest{
		UserId:        m.Author.ID,
		TransactionId: txid,
		AmountInUsdt:  20,
		AuditNote:     paymentsproto.PaymentTypeFuturesSubscription,
	}).Send(ctx).Response()
	switch {
	case gerrors.Is(err, gerrors.ErrFailedPrecondition, "failed_to_register_payment.user_does_not_have_an_account"):
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			":disappointed: Hey, you must have an account registered first! You can call `!help` to see how.",
		)
		return err
	case gerrors.Is(err, gerrors.ErrAlreadyExists, "failed_to_register_payment.payment_already_exists"):
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			":rocket: Hey, looks like you've already paid for this month! If you don't think that's correct please ping @ajperkins",
		)
		return err
	case gerrors.Is(err, gerrors.ErrFailedPrecondition, "failed_to_register_payment.user_has_already_paid"):
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			":rocket: Hey, looks like you've already paid for this month! If you don't think that's correct please ping @ajperkins",
		)
		return err
	case gerrors.Is(err, gerrors.ErrFailedPrecondition, "failed_to_register_payment.transaction_of_correct_amount_does_not_exist_in_deposit_account"):
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			":disappointed: Hey, my bad I can't find that transaction id in the deposit account! Please check that the **transaction id** and the **amount** is correct ",
		)
		return err
	case err != nil:
		slog.Error(ctx, "Failed to process payment: %v", err.Error())
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			":disappointed: Hey, apologies! Looks like something broke. Please try again - if this keeps happening please ping @ajperkins to investigate",
		)
		return err
	}

	_, err = s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: Payment successfully registered. Thank you <@%s>! :coin:", m.Author.ID),
	)

	return err
}
