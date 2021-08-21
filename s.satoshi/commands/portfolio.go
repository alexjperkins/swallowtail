package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	googlesheetsproto "swallowtail/s.googlesheets/proto"
)

const (
	portfolioCommandID     = "porfolio-command"
	portfolioCommandPrefix = "!portfolio"
	portfolioCommandUsage  = `
	Usage: !portfolio <operation> <args>
	Example: !portfolio create alexperkins.crypto@gmail.com

	Operations:
	1. create: creates a new portfolio in googlesheets: args [email_address].
	2. list: list all your googlesheets: args [].
	3. register: register a previously created sheet: args [googlesheet_url, sheet_name, email_address]
	4. delete: deletes a googlesheet by googlesheet id, you can find this by the "list" command: args [googlesheet_id]
	`
)

func init() {
	register(portfolioCommandID, &Command{})
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

	switch strings.ToLower(tokens[1]) {
	case "create":
		if len(tokens) < 3 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@:%s>, bad args! Usage: %s", m.Author.ID, portfolioCommandUsage))
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
	case "register":
		if len(tokens) < 5 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@:%s>, bad args! Usage: %s", m.Author.ID, portfolioCommandUsage))
		}

		url, sheetName, email := tokens[2], tokens[3], tokens[4]
		shareEmail, err := registerPortfolioSheet(m.Author.ID, url, sheetName, email)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				":wave: <@:%s>, Failed to register sheet! Please check that the url is correct", m.Author.ID,
			))
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			":rocket: <@:%s>, Googlesheet registered! Please share the sheet with this email to allow satoshi access to sync: %s", m.Author.ID, shareEmail,
		))

	case "delete":
		if len(tokens) < 3 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@:%s>, bad args! Usage: %s", m.Author.ID, portfolioCommandUsage))
		}

		sheetID := tokens[2]
		if err := deletePortfolioSheet(sheetID, m.Author.ID); err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				":wave: <@:%s>, Failed to delete sheet! Please check that the id is correct and the sheet exists", m.Author.ID,
			))
		}
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

func registerPortfolioSheet(userID, url, sheetName, email string) (string, error) {
	rsp, err := (&googlesheetsproto.RegisterNewPortfolioSheetRequest{
		UserId:    userID,
		Url:       url,
		SheetName: sheetName,
		Email:     email,
	}).Send(context.Background()).Response()
	if err != nil {
		return "", err
	}

	return rsp.GetServiceAccountEmail(), nil
}

func deletePortfolioSheet(googlesheetID, userID string) error {
	if _, err := (&googlesheetsproto.DeleteSheetBySheetIDRequest{
		GooglesheetId: googlesheetID,
		UserId:        userID,
	}).Send(context.Background()).Response(); err != nil {
		return err
	}
	return nil
}
