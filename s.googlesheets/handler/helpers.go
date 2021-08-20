package handler

import (
	"regexp"
	"swallowtail/libraries/gerrors"
)

var (
	re = regexp.MustCompile(`^https://docs.google.com/spreadsheets/d/([a-zA-Z0-9-_]{44})/edit[a-zA-Z0-9#=]*$`)
)

func parseSpreadsheetIDFromURL(url string) (string, error) {
	matches := re.FindStringSubmatch(url)
	if len(matches) != 2 {
		return "", gerrors.BadParam("invalid_google_spreadsheets_url", map[string]string{
			"spreadsheet_url": url,
		})
	}

	return matches[1], nil
}
