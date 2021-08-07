package marshaling

import (
	"google.golang.org/api/sheets/v4"

	"swallowtail/s.googlesheets/domain"
	googlesheetsproto "swallowtail/s.googlesheets/proto"
)

// CreatePortfolioProtoToDomain marshals the incoming proto to a googlesheets
func CreatePortfolioProtoToDomain(in *googlesheetsproto.CreatePortfolioSheetRequest, ss *sheets.Spreadsheet) *domain.Googlesheet {
	var sheetID = "SpotPortfolio"
	if len(ss.Sheets) > 0 {
		sheetID = ss.Sheets[0].Properties.Title
	}
	return &domain.Googlesheet{
		SheetID:           sheetID,
		Email:             in.Email,
		SpreadsheetID:     ss.SpreadsheetId,
		URL:               ss.SpreadsheetUrl,
		UserID:            in.UserId,
		Active:            in.Active,
		SheetType:         "PORTFOLIO",
		WithPagerOnError:  in.ShouldPagerOnError,
		WithPagerOnTarget: in.ShouldPagerOnTarget,
	}
}
