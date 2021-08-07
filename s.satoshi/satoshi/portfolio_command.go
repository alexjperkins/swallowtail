package satoshi

import (
	"context"
	"fmt"
	"strings"
	"time"

	googlesheetsproto "swallowtail/s.googlesheets/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	portfolioCommandID     = "porfolio-command"
	portfolioCommandPrefix = "!portfolio"
	portfolioCommandUsage  = `
	Usage: !portfolio <operation> <args>
	Example !portfolio create alexperkins.crypto@gmail.com

	Operations:
	1: create: creates a new portfolio in googlesheets: args [email_address].
	2: list: list all your googlesheets: args [].
	`
)

func init() {
	registerSatoshiCommand(portfolioCommandID, portfolioCommand)
}

func portfolioCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, portfolioCommandPrefix) {
		return
	}

	tokens := strings.Split(m.Content, " ")
	if len(tokens) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@:%s, Usage: %s>", m.Author.ID, portfolioCommandUsage))
		return
	}
	slog.Info(context.TODO(), "Received %s command, args: %v", priceCommandPrefix, tokens)

	switch strings.ToLower(tokens[1]) {
	case "create":
		if len(tokens) < 3 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@:%s, bad args! Usage: %s>", m.Author.ID, portfolioCommandUsage))
			return
		}

		email := tokens[2]
		url, err := createPortfolioSheet(m.Author.ID, email)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s>, I failed to create a portfolio sheet: %v", m.Author.ID, err))
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":rocket: <@%s>, Portfolio sheet created, here's the URL: %s", m.Author.ID, url))
		return
	case "list":
		sheets, err := listPortfolioSheets(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":disappointed: <@%s>, I failed to list your portfolio sheets: %v", m.Author.ID, err))
			return
		}

		sheetsMsgArr := []string{"URL: Type"}
		for _, sheet := range sheets {
			// Todo Add padding
			sheetsMsgArr = append(sheetsMsgArr, fmt.Sprintf("%s: %s", sheet.Url, sheet.SheetType))
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s>, here are your sheets\n\n%s", m.Author.ID, strings.Join(sheetsMsgArr, "\n")))
	}
}

func createPortfolioSheet(userID, email string) (string, error) {
	rsp, err := (&googlesheetsproto.CreatePortfolioSheetRequest{
		UserId:              userID,
		Email:               email,
		Active:              true,
		ShouldPagerOnError:  true,
		ShouldPagerOnTarget: true,
	}).SendWithTimeout(context.Background(), 15*time.Second).Response()
	if err != nil {
		return "", err
	}

	return rsp.GetURL(), nil
}

func listPortfolioSheets(userID string) ([]*googlesheetsproto.SheetResponse, error) {
	rsp, err := (&googlesheetsproto.ListSheetsByUserIDRequest{
		UserId: userID,
	}).Send(context.Background()).Response()
	if err != nil {
		return nil, err
	}
	return rsp.Sheets, nil
}
