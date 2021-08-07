package templates

// SheetType  ...
type SheetType int

const (
	// PortfolioSheetType defines the protofolio sheet.
	PortfolioSheetType SheetType = iota
	UnknownSheetType
	NoTemplateType
)

func (s SheetType) String() string {
	switch s {
	case PortfolioSheetType:
		return "PORTFOLIO"
	default:
		return "UNKNOWN"
	}
}

// GooglesheetsTemplate defines the interface of a googlesheets template
type GooglesheetsTemplate interface {
	ID() SheetType
	RowRange() string
	Values() [][]interface{}
}
