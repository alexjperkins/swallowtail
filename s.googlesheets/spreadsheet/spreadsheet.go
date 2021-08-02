package spreadsheet

import (
	"context"
)

var (
	stringDatatype  = "string"
	float64Datatype = "float64"
	intDatatype     = "int"
)

type GoogleSpreadsheet interface {
	CreateNew(sheetName string) (string, error)
	Rows(context.Context, string, string) []interface{}
	ParseRow(context.Context, []interface{})
	UpdateRow(context.Context, interface{})
	UpdateRows(context.Context, []interface{})
}
