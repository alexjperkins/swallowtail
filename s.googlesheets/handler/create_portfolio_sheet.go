package handler

import (
	"context"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	"swallowtail/s.googlesheets/client"
	"swallowtail/s.googlesheets/dao"
	"swallowtail/s.googlesheets/marshaling"
	googlesheetsproto "swallowtail/s.googlesheets/proto"
	"swallowtail/s.googlesheets/templates"
)

// CreatePortfolioSheet creates a new portfolio sheet & registers to sync.
func (s *GooglesheetsService) CreatePortfolioSheet(
	ctx context.Context, in *googlesheetsproto.CreatePortfolioSheetRequest,
) (*googlesheetsproto.CreatePortfolioSheetResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.Email == "":
		return nil, gerrors.BadParam("missing_param.email", nil)
	}

	errParams := map[string]string{
		"user_id":       in.UserId,
		"email_address": in.Email,
	}

	// Check the user first has an account registered.
	_, err := (&accountproto.ReadAccountRequest{
		UserId: in.UserId,
	}).Send(ctx).Response()
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		return nil, gerrors.FailedPrecondition("account_not_registered: User cannot create portfolio without registered account", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "Failed to create portfolio", errParams)
	}

	// Prevent user spam
	sheets, err := dao.ListSheetsByUserID(ctx, in.GetUserId())
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "not_found.no_googlesheets_registered_for_user"):
		// This is fine.
	case err != nil:
		return nil, gerrors.Augment(err, "Failed to create portfolio sheet; couldn't check existing sheets for user", errParams)
	}

	if len(sheets) >= 5 {
		return nil, gerrors.FailedPrecondition("max_portfolio_sheets_reached", errParams)
	}

	// Create our portfolio sheet.
	ss, err := client.CreateSheet(ctx, templates.PortfolioSheetType, in.GetEmail())
	if err != nil {
		return nil, gerrors.Augment(err, "Failed to create googlesheet", errParams)
	}

	slog.Info(ctx, "Created sheet: %s, %s", in.Email, ss.SpreadsheetUrl)

	// Lets persist the googlesheets now that we've created one. This will allow us to sync later.
	sheet := marshaling.CreatePortfolioProtoToDomain(in, ss)

	slog.Warn(nil, "%+v", sheet)

	if err := dao.RegisterGooglesheet(ctx, sheet); err != nil {
		return nil, gerrors.Augment(err, "Failed to create googlesheet", errParams)
	}

	return &googlesheetsproto.CreatePortfolioSheetResponse{
		URL: sheet.URL,
	}, nil
}
