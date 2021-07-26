package templates

func init() {
	registerTemplate(PortfolioSheetType, portfolioTemplate{})
}

type portfolioTemplate struct{}

func (p portfolioTemplate) ID() SheetType {
	return PortfolioSheetType
}

func (p portfolioTemplate) RowRange() string {
	return ""
}

func (p portfolioTemplate) Values() [][]interface{} {
	return [][]interface{}{}
}
