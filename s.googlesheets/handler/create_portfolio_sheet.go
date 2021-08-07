package handler

import (
	"context"

	"github.com/monzo/terrors"

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
		return nil, terrors.PreconditionFailed("missing-param.user-id", "Failed to create googlesheet; missing user id", nil)
	case in.ShouldShare && in.Email == "":
		return nil, terrors.BadRequest("missing-param.email", "Failed to create googlesheet; cannot share with no email", nil)
	}

	errParams := map[string]string{
		"user_id": in.UserId,
	}

	// Prevent user spam
	sheets, err := dao.ListSheetsByUserID(ctx, in.GetUserId())
	if err != nil {
		return nil, terrors.Augment(err, "Failed to create portfolio sheet; couldn't check existing sheets for user", errParams)
	}
	if len(sheets) >= 5 {
		return nil, terrors.PreconditionFailed("max-portfolio-sheets-reached", "User already has a maximum of 5 portfolio sheets", errParams)
	}

	// Create our portfolio sheet.
	url, err := client.CreateSheet(ctx, templates.PortfolioSheetType)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to create googlesheet", errParams)
	}

	// Lets persist the googlesheets now that we've created one. This will allow us to sync later.
	sheet := marshaling.CreatePortfolioProtoToDomain(in)
	if err := dao.RegisterGooglesheet(ctx, sheet); err != nil {
		return nil, terrors.Augment(err, "Failed to create googlesheet", errParams)
	}

	return &googlesheetsproto.CreatePortfolioSheetResponse{
		URL: url,
	}, nil
}
