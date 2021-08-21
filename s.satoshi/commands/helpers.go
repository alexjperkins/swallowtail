package commands

import (
	"fmt"
	"strings"
	"swallowtail/libraries/util"
	accountproto "swallowtail/s.account/proto"

	"github.com/bwmarrin/discordgo"
)

func formatExchangesToMsg(exchanges []*accountproto.Exchange, m *discordgo.MessageCreate) string {
	var lines = []string{}
	lines = append(lines, "`Exchange: ID Username MaskedAPIKey MaskedSecretKey`")
	for i, exchange := range exchanges {
		// We're masking here to be on the safe side; we should expect them to already be masked.
		// TODO maybe we should ping someone here or something.
		maskedAPIKey, maskedSecretKey := util.MaskKey(exchange.ApiKey, 4), util.MaskKey(exchange.SecretKey, 4)

		line := fmt.Sprintf(
			"`%v) %s: %s %s %s %s`",
			i, exchange.Exchange, exchange.ExchangeId, m.Author.Username, maskedAPIKey, maskedSecretKey,
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
