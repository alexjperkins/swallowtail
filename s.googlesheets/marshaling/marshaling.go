package marshaling

import (
	"swallowtail/s.googlesheets/domain"
	googlesheetsproto "swallowtail/s.googlesheets/proto"
)

// CreatePortfolioProtoToDomain marshals the incoming proto to a googlesheets
func CreatePortfolioProtoToDomain(in *googlesheetsproto.CreatePortfolioSheetRequest) *domain.Googlesheet {
	return &domain.Googlesheet{
		UserID:            in.UserId,
		Active:            in.Active,
		WithPagerOnError:  in.ShouldPagerOnError,
		WithPagerOnTarget: in.ShouldPagerOnTarget,
	}
}
