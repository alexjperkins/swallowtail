package handler

import (
	googlesheetsproto "swallowtail/s.googlesheets/proto"
)

// GooglesheetsService ...
type GooglesheetsService struct {
	*googlesheetsproto.UnimplementedGooglesheetsServer
}
