package consumers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/monzo/slog"
)

func formatContent(ctx context.Context, username, timestamp, content string) string {
	// Parse timestamp.
	ts, err := time.Parse(time.RFC3339, timestamp)
	switch {
	case err != nil:
		slog.Warn(ctx, "Failed to parse timestamp; setting as original: %s, err: %v", timestamp, err)
	default:
		timestamp = ts.Truncate(time.Minute).String()
	}

	header := fmt.Sprintf(":dove: `MOD MESSAGE:` %s   :brain:", username)
	tpl := `
Timestamp:    %v

Content:
%s
`
	attachmentsTpl := `
Attachments:
`

	// Parse attachments.
	splits := strings.Split(content, "[Attachments]")
	var attachments string
	if len(splits) > 1 {
		content, attachments = splits[0], strings.Join(splits[1:], " ")
	}

	formattedContent := fmt.Sprintf(tpl, timestamp, content)
	formattedAttachments := fmt.Sprintf("```%s```%s", attachmentsTpl, attachments)

	if attachments == "" {
		return fmt.Sprintf("%s```%s```", header, formattedContent)
	}

	return fmt.Sprintf("%s```%s```%s", header, formattedContent, formattedAttachments)
}
