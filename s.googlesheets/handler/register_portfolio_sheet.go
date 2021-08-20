package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	"swallowtail/s.googlesheets/dao"
	"swallowtail/s.googlesheets/domain"
	googlesheetsproto "swallowtail/s.googlesheets/proto"
	"swallowtail/s.googlesheets/templates"
)

// RegisterNewPortfolioSheet ...
func (s *GooglesheetsService) RegisterNewPortfolioSheet(
	ctx context.Context, in *googlesheetsproto.RegisterNewPortfolioSheetRequest,
) (*googlesheetsproto.RegisterNewPortfolioSheetResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.FailedPrecondition("missing_param.user_id", nil)
	case in.SheetName == "":
		return nil, gerrors.FailedPrecondition("missing_param.sheet_name", nil)
	case in.Url == "":
		return nil, gerrors.FailedPrecondition("missing_param.spreadsheet_id", nil)
	case in.Email == "":
		return nil, gerrors.FailedPrecondition("missing_param.email", nil)
	}

	spreadsheetID, err := parseSpreadsheetIDFromURL(in.Url)
	if err != nil {
		return nil, gerrors.Augment(err, "parsing_failed.spreadsheet_id", nil)
	}

	errParams := map[string]string{
		"user_id":        in.UserId,
		"spreadsheet_id": spreadsheetID,
		"sheet_name":     in.SheetName,
	}

	// Check the user first has an account registered.
	_, err = (&accountproto.ReadAccountRequest{
		UserId: in.UserId,
	}).Send(ctx).Response()
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account-not-found"):
		return nil, gerrors.FailedPrecondition("account-not-registered: User cannot create portfolio without registered account", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "Failed to create portfolio", errParams)
	}

	// List the number of sheets the user already has; if they already have 5, we exit since they have
	// already reached the limit.
	sheets, err := dao.ListSheetsByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "not_found.no-googlesheets-registered-for-user"):
		// This is fine.
	case err != nil:
		return nil, gerrors.Augment(err, "Failed to create portfolio sheet; couldn't check existing sheets for user", errParams)
	}

	if len(sheets) >= 15 {
		return nil, gerrors.FailedPrecondition("maximum_sheets_reached_for_user", errParams)
	}

	if err := (dao.RegisterGooglesheet(ctx, &domain.Googlesheet{
		UserID:            in.UserId,
		SpreadsheetID:     spreadsheetID,
		SheetID:           in.SheetName,
		SheetType:         templates.PortfolioSheetType.String(),
		Active:            true,
		WithPagerOnError:  true,
		WithPagerOnTarget: true,
		Email:             in.Email,
		URL:               in.Url,
	})); err != nil {
		return nil, gerrors.Augment(err, "failed_to_register_googlesheet", errParams)
	}

	return &googlesheetsproto.RegisterNewPortfolioSheetResponse{}, nil
}
