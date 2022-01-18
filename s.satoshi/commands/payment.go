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
	paymentUsage     = `!payment <subcommand>`
)

func init() {
	register(paymentCommandID, &Command{
		ID:                  paymentCommandID,
		IsPrivate:           true,
		MinimumNumberOfArgs: 1,
		Usage:               paymentUsage,
		Handler:             paymentHandler,
		Description:         "Command for registering subscriptions & viewing payments.",
		Guide:               "https://scalloped-single-1bd.notion.site/How-to-register-a-subscription-payment-35abb69004d946de8010c2f58d9863e1",
		SubCommands: map[string]*Command{
			"register": {
				ID:                  "payment-register",
				IsPrivate:           true,
				MinimumNumberOfArgs: 1,
				Usage:               `!payment register <transaction_id>`,
				Description:         "Registers a new payment to satoshi. It checks the transaction is correct & keeps a record of it.",
				Handler:             registerPaymentHandler,
				FailureMsg:          "Please check that you have an account registered; or maybe you've already paid for this month? ping @ajperkins if unsure",
			},
			"uptodate": {
				ID:                  "payment-up-to-date",
				IsPrivate:           true,
				MinimumNumberOfArgs: 0,
				Usage:               `!payment up-to-date`,
				Description:         "Checks if you are up to date on payments with the bot.",
				Handler:             upToDatePaymentHandler,
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

func upToDatePaymentHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	var msg string
	rsp, err := (&paymentsproto.ReadUsersLastPaymentRequest{
		UserId:  m.Author.ID,
		ActorId: paymentsproto.ActorSatoshiSystem,
	}).Send(ctx).Response()
	switch {
	case gerrors.PartialIs(err, gerrors.ErrUnknown, "no_rows_in_result_set"):
		msg = "I can't find any payment I'm afraid :grimacing:"
	case err != nil:
		return gerrors.Augment(err, "failed_up_to_date_payment_command", nil)
	case rsp.GetHasUserPaidForLastMonth():
		msg = "it looks like you've already paid for the month! :rocket: :dove:"
	default:
		msg = "I can't find a payment for the last month I'm afraid. Please ask in support channels if you think this is wrong!"
	}

	if _, err := s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: <@%s> %s", m.Author.ID, msg),
	); err != nil {
		slog.Error(ctx, "Failed to publish up to date message to user via discord: %s", m.Author.Username)
	}

	return nil
}
