package commands

import (
	"fmt"
	"strconv"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"

	"github.com/bwmarrin/discordgo"
)

func formatUsageMsg(userID, usage string, guide string) string {
	var formattedGuide string
	if guide != "" {
		formattedGuide = fmt.Sprintf("Guide: `%s`", guide)
	}

	return fmt.Sprintf(":wave: <@%s>, that's not quite how the command works.\n%s\n%s", userID, formattedGuide, util.WrapAsCodeBlock(usage))
}

func formatNonAdminMsg(userID string) string {
	return fmt.Sprintf(":wave: <@%s>, apologies! But this command can only be run by admins :disappointed:", userID)
}

func formatNonFuturesMsg(userID string) string {
	return fmt.Sprintf(":wave: <@%s>, apologies! But this command can only be run by futures members :grimacing: Ping @ajperkins if you want to know how to become one", userID)
}

func formatNonPublicMsg(userID string) string {
	return fmt.Sprintf(":wave: <@%s>, Please DM satoshi this command instead, the response may contain sensitive information. Thanks", userID)
}

func formatFailureMsg(userID, usage, failureMsg string, err error) string {
	var errMsg = err.Error()
	switch {
	case gerrors.Is(err, gerrors.ErrUnimplemented):
		errMsg = "Command unimplemented"
	}

	return fmt.Sprintf(
		":disappointed: Sorry <@%s>, I failed to execute that command.\n%s\n Error: %s\n",
		userID, failureMsg, errMsg,
	)
}

func formatHelpMsg(command *Command, isFuturesMember, isAdmin bool) string {
	switch {
	case command.IsFuturesOnly && !isFuturesMember:
		return ""
	case command.IsAdminOnly && !isAdmin:
		return ""
	case command.Usage == "":
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n\n%s", strings.ToUpper(command.ID)))
	sb.WriteString(fmt.Sprintf("\n\tUsage: %s", command.Usage))

	if command.Description != "" {
		sb.WriteString(fmt.Sprintf("\n\tDescription: %s", command.Description))
	}

	if command.Guide != "" {
		sb.WriteString(fmt.Sprintf("\n\tGuide: %s", command.Guide))
	}

	sb.WriteString("\n")

	if len(command.SubCommands) > 0 {
		sb.WriteString("\n\tSubcommands")
		for id, sc := range command.SubCommands {
			sb.WriteString(fmt.Sprintf("\n\t\t%s: %s", id, sc.Description))
		}
	}

	return sb.String()
}

func formatExchangesToMsg(exchanges []*accountproto.Exchange, m *discordgo.MessageCreate) string {
	var lines = []string{}
	lines = append(lines, "`Exchange: ID Username MaskedAPIKey MaskedSecretKey`")
	for i, exchange := range exchanges {
		// We're masking here to be on the safe side; we should expect them to already be masked.
		// TODO maybe we should ping someone here or something.
		maskedAPIKey, maskedSecretKey := util.MaskKey(exchange.ApiKey, 4), util.MaskKey(exchange.SecretKey, 4)

		line := fmt.Sprintf(
			"`%v) %s: %s %s %s %s`",
			i, exchange.ExchangeType, exchange.ExchangeId, m.Author.Username, maskedAPIKey, maskedSecretKey,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// calculateRisk returns the number of contracts to buy.
func calculateRisk(entry, stopLoss, accountSize, percentage float64) float64 {
	switch {
	case entry == stopLoss:
		return 0.0
	}

	maxRiskToLose := percentage * accountSize
	lossPerContract := entry - stopLoss

	return maxRiskToLose / lossPerContract
}

func contains(needle string, haystack []string) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}
	return false
}

func safeTokenParse(tokens []string, index int) (string, error) {
	if index >= len(tokens) {
		return "", gerrors.FailedPrecondition("failed_safe_token_parse", map[string]string{
			"index":      strconv.Itoa(index),
			"len_tokens": strconv.Itoa(len(tokens)),
		})
	}

	return tokens[index], nil
}

func isFuturesMember(roles []string) bool {
	for _, role := range roles {
		if role == discordproto.DiscordSatoshiFuturesRoleID {
			return true
		}
	}

	return false
}

func isAdmin(roles []string) bool {
	for _, role := range roles {
		if role == discordproto.DiscordSatoshiAdminRoleID {
			return true
		}
	}

	return false
}

// Placeholder
func normalizeContent(content string) string {
	return content
}

func getMembersRolesFromGuild(session *discordgo.Session, userID string) ([]string, error) {
	m, err := session.GuildMember(discordproto.DiscordSatoshiGuildID, userID)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_members_guild_roles", nil)
	}

	return m.Roles, nil
}
