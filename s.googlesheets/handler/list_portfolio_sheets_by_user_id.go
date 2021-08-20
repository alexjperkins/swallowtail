package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.googlesheets/dao"
	googlesheetsproto "swallowtail/s.googlesheets/proto"
)

// ListSheetsByUserID ...
func (s *GooglesheetsService) ListSheetsByUserID(
	ctx context.Context, in *googlesheetsproto.ListSheetsByUserIDRequest,
) (*googlesheetsproto.ListSheetsByUserIDResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.FailedPrecondition("missing_param.user_id", nil)
	}

	errParams := map[string]string{
		"user_id": in.UserId,
	}

	sheets, err := dao.ListSheetsByUserID(ctx, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "Failed to list sheets by user id", errParams)
	}

	protoSheets := []*googlesheetsproto.SheetResponse{}
	for _, sheet := range sheets {
		protoSheets = append(protoSheets, &googlesheetsproto.SheetResponse{
			Url:       sheet.URL,
			SheetType: sheet.SheetType,
			SheetName: sheet.SheetID,
			SheetId:   sheet.GooglesheetID,
		})
	}

	return &googlesheetsproto.ListSheetsByUserIDResponse{
		Sheets: protoSheets,
	}, nil
}
