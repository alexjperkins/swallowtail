package pricebot

import (
	"fmt"
	"strings"
	"swallowtail/libraries/util"
	"time"
)

var (
	lineTemplate     = "[%s]: [%s]"
	greetingTemplate = "\n:robot: **Price bot hourly update** :robot:\n[%v]\n\nPlease ping **@ajperkins** if you'd like a coin adding.\n"
)

func buildMessage(prices []*PriceBotPrice, withGreeting bool) string {
	if len(prices) == 0 {
		return ""
	}
	lines := []string{}
	for _, price := range prices {
		line := buildLine(price)
		lines = append(lines, line)
	}

	strLines := strings.Join(lines, "\n")
	if !withGreeting {
		return monospaceWrapper(strLines)
	}
	return fmt.Sprintf("%s%s", buildGreeting(), monospaceWrapper(strLines))
}

func buildLine(price *PriceBotPrice) string {
	strPrice, _ := util.FormatPriceAsString(price.Price)
	if price.Price == 0.0 {
		strPrice = "N/A"
	}
	return fmt.Sprintf(lineTemplate, price.Symbol, strPrice)
}

func buildGreeting() string {
	return fmt.Sprintf(greetingTemplate, time.Now().Round(time.Minute))
}

func monospaceWrapper(s string) string {
	return fmt.Sprintf("```%s```", s)
}
