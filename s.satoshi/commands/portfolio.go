package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"swallowtail/libraries/gerrors"
	googlesheetsproto "swallowtail/s.googlesheets/proto"
)

const (
	portfolioCommandID    = "porfolio"
	portfolioCommandUsage = `
	Usage: !portfolio <subcommand> <args>
	Example: !portfolio create alexperkins.crypto@gmail.com

	Subcommands:
	1. create: creates a new portfolio in googlesheets: args [email_address].
	2. list: list all your googlesheets: args [].
	3. register: register a previously created sheet: args [googlesheet_url, sheet_name, email_address]
	4. delete: deletes a googlesheet by googlesheet id, you can find this by the "list" command: args [googlesheet_id]
	`
)

func init() {
	register(portfolioCommandID, &Command{
		ID:                  portfolioCommandID,
		Private:             true,
		MinimumNumberOfArgs: 1,
		FailureMsg:          "",
		Handler:             portfolioCommand,
		Usage:               portfolioCommandUsage,
		SubCommands: map[string]*Command{
			"list": {
				ID:         "portfolio-list",
				Private:    true,
				Usage:      `!portfolio list`,
				FailureMsg: "",
				Handler:    listPortfolioCommand,
			},
			"delete": {
				ID:                  "portfolio-delete",
				Private:             true,
				Usage:               "!portfolio delete <googlesheet_id>",
				FailureMsg:          "Please check the googlesheet id is correct - you can find this by using the list command",
				Handler:             deletePorfolioCommand,
				MinimumNumberOfArgs: 1,
			},
			"create": {
				ID:                  "portfolio-create",
				Private:             true,
				Usage:               "!portfolio create <email>",
				MinimumNumberOfArgs: 1,
				FailureMsg:          "Please check your email is correct; if it is ping @ajperkins with the error message",
				Handler:             createPortfolioCommand,
			},
			"register": {
				ID:                  "portfolio-register",
				Private:             true,
				Usage:               "!portfolio register <url> <sheet_name> <email>",
				MinimumNumberOfArgs: 3,
				FailureMsg:          "Please check the URL & sheet name passed are correct; satoshi needs the full URL.",
				Handler:             registerPortfolioCommand,
			},
		},
	})
}

func portfolioCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return gerrors.Unimplemented("parent_command_unimplemented.portfolio", nil)
}

func listPortfolioCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	rsp, err := (&googlesheetsproto.ListSheetsByUserIDRequest{
		UserId: m.Author.ID,
	}).Send(context.Background()).Response()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":disappointed: <@%s>, I failed to list your portfolio sheets: %v", m.Author.ID, err))
		return gerrors.Augment(err, "failed_to_list_sheets", map[string]string{
			"discord_username": m.Author.Username,
			"user_id":          m.Author.ID,
		})
	}

	sheetsMsgArr := []string{"URL: Type"}
	for _, sheet := range rsp.GetSheets() {
		// Todo Add padding
		sheetsMsgArr = append(sheetsMsgArr, fmt.Sprintf("%s: %s", sheet.Url, sheet.SheetType))
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s>, here are your sheets\n\n%s", m.Author.ID, strings.Join(sheetsMsgArr, "\n")))
	return nil
}

func deletePorfolioCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	googlesheetID := tokens[0]
	if _, err := (&googlesheetsproto.DeleteSheetBySheetIDRequest{
		GooglesheetId: googlesheetID,
		UserId:        m.Author.ID,
	}).Send(context.Background()).Response(); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			":wave: <@:%s>, Failed to delete sheet! Please check that the id is correct and the sheet exists", m.Author.ID,
		))
		return gerrors.Augment(err, "failed_to_delete_sheet", map[string]string{
			"discord_username": m.Author.Username,
			"user_id":          m.Author.ID,
			"googlesheet_id":   googlesheetID,
		})
	}

	return nil
}

func createPortfolioCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	email := tokens[0]
	rsp, err := (&googlesheetsproto.CreatePortfolioSheetRequest{
		UserId:              m.Author.ID,
		Email:               email,
		Active:              true,
		ShouldPagerOnError:  true,
		ShouldPagerOnTarget: true,
	}).SendWithTimeout(context.Background(), 15*time.Second).Response()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s>, I failed to create a portfolio sheet: %v", m.Author.ID, err))
		return gerrors.Augment(err, "failed_to_create_sheet", map[string]string{
			"discord_username": m.Author.Username,
			"user_id":          m.Author.ID,
			"email":            email,
		})
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":rocket: <@%s>, Portfolio sheet created, here's the URL: %s", m.Author.ID, rsp.GetURL()))

	return nil
}

func registerPortfolioCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	url, sheetName, email := tokens[0], tokens[1], tokens[2]
	rsp, err := (&googlesheetsproto.RegisterNewPortfolioSheetRequest{
		UserId:    m.Author.ID,
		Url:       url,
		SheetName: sheetName,
		Email:     email,
	}).Send(context.Background()).Response()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			":wave: <@:%s>, Failed to register sheet! Please check that the url is correct", m.Author.ID,
		))
		return gerrors.Augment(err, "failed_to_register_sheet", map[string]string{
			"discord_username": m.Author.Username,
			"user_id":          m.Author.ID,
			"url":              url,
			"sheet_name":       sheetName,
		})
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
		":rocket: <@:%s>, Googlesheet registered! Please share the sheet with this email to allow satoshi access to sync: %s", m.Author.ID, rsp.GetServiceAccountEmail(),
	))

	return nil
}
