package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
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
	case in.SpreadsheetId == "":
		return nil, gerrors.FailedPrecondition("missing_param.spreadsheet_id", nil)
	}

	errParams := map[string]string{
		"user_id":        in.UserId,
		"spreadsheet_id": in.SpreadsheetId,
		"sheet_name":     in.SheetName,
	}

	sheets, err := dao.ListSheetsByUserID(ctx, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_users_sheets", errParams)
	}

	if len(sheets) >= 5 {
		return nil, gerrors.FailedPrecondition("maximum-sheets-reached-for-user", errParams)
	}

	if err := (dao.RegisterGooglesheet(ctx, &domain.Googlesheet{
		UserID:            in.UserId,
		SpreadsheetID:     in.SpreadsheetId,
		SheetID:           in.SheetName,
		SheetType:         templates.PortfolioSheetType.String(),
		Active:            true,
		WithPagerOnError:  true,
		WithPagerOnTarget: true,
		Email:             in.Email,
		URL:               "",
	})); err != nil {
		return nil, gerrors.Augment(err, "failed-to-register-googlesheet", errParams)
	}

	return &googlesheetsproto.RegisterNewPortfolioSheetResponse{}, nil
}
