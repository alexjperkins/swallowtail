package spreadsheet

import (
	"context"
)

var (
	stringDatatype  = "string"
	float64Datatype = "float64"
	intDatatype     = "int"
)

type GoogleSpreadsheet struct {
	// The ID of the spreadsheet itself
	ID string
	// The list of self contained sheet IDs
	SheetIDs []GoogleSheet
}

type GoogleSheet interface {
	ID() string
	CreateNew(sheetName string) error
	Rows(context.Context, string, string) []interface{}
	ParseRow(context.Context, []interface{})
	UpdateRow(context.Context, interface{})
	UpdateRows(context.Context, []interface{})
}
