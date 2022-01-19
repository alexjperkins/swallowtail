package commands

import (
	"context"
	"fmt"
	"strings"
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
			"check": {
				ID:                  "payment-check",
				IsPrivate:           true,
				MinimumNumberOfArgs: 0,
				Usage:               `!payment check`,
				Description:         "Checks if you are up to date on payments with the bot.",
				Handler:             checkPaymentHandler,
			},
			"list": {
				ID:                  "payment-list",
				IsPrivate:           true,
				MinimumNumberOfArgs: 0,
				Usage:               `!payment list`,
				Description:         "Lists all payments registered with the bot for the last 10 years",
				Handler:             listPaymentHandler,
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

func checkPaymentHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
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
		msg = fmt.Sprintf("it looks like you've already paid for the month, at %s! :rocket: :dove:", rsp.GetLastPaymentTimestamp().AsTime())
	default:
		msg = fmt.Sprintf("I can't seem to find a payment for the last month I'm afraid, your last was: `%s`.\nPlease ask in support channels if you think this is wrong!", rsp.GetLastPaymentTimestamp().AsTime())
	}

	if _, err := s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: <@%s> %s", m.Author.ID, msg),
	); err != nil {
		slog.Error(ctx, "Failed to publish up to date message to user via discord: %s", m.Author.Username)
	}

	return nil
}

func listPaymentHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	rsp, err := (&paymentsproto.ListPaymentsByUserIDRequest{
		UserId:  m.Author.ID,
		ActorId: paymentsproto.ActorSatoshiSystem,
		Limit:   24, // This gives us the last 2 years of payments.
	}).Send(ctx).Response()
	if err != nil {
		return gerrors.Augment(err, "failed_to_list_payments", nil)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nTotal Amount (USDT): %.2f\n", sumPayments(rsp.GetPayments())))
	sb.WriteString(fmt.Sprintf("Number of payments:  %d\n", len(rsp.GetPayments())))

	sb.WriteString("\nPayments:\n")

	for _, payment := range rsp.GetPayments() {
		sb.WriteString(fmt.Sprintf("\n[%s] txid: %s amount: %.2f", payment.GetPaymentTimestamp().AsTime(), payment.GetTransactionId(), payment.GetAmountInUsdt()))
	}

	if _, err := s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: <@%s> I've found your payments over the last two years :dove:\n```%s```", m.Author.ID, sb.String()),
	); err != nil {
		slog.Error(ctx, "Failed to send user [%s] list payments message", m.Author.Username)
	}

	return nil
}

func sumPayments(payments []*paymentsproto.Payment) float32 {
	if len(payments) == 0 {
		return 0
	}

	return payments[0].GetAmountInUsdt() + sumPayments(payments[1:])
}
