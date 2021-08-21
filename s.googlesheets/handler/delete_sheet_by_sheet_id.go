package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.googlesheets/dao"
	googlesheetsproto "swallowtail/s.googlesheets/proto"

	"github.com/monzo/slog"
)

// DeleteSheetBySheetID ...
func (s *GooglesheetsService) DeleteSheetBySheetID(
	ctx context.Context, in *googlesheetsproto.DeleteSheetBySheetIDRequest,
) (*googlesheetsproto.DeleteSheetBySheetIDResponse, error) {
	switch {
	case in.GooglesheetId == "":
		return nil, gerrors.BadParam("missing_param.spreadsheet_id", nil)
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	}

	if err := dao.DeleteSheetBySheetID(ctx, in.UserId, in.GooglesheetId); err != nil {
		return nil, gerrors.Augment(err, "delete_sheet_failed", map[string]string{
			"user_id":        in.UserId,
			"googlesheet_id": in.GooglesheetId,
		})
	}

	slog.Info(ctx, "Deleted sheet: %s user: %s", in.GooglesheetId, in.UserId)

	return &googlesheetsproto.DeleteSheetBySheetIDResponse{}, nil
}
