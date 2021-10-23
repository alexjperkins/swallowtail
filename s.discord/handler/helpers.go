package handler

import (
	"fmt"
	"io"
	"strings"
	"swallowtail/libraries/gerrors"
)

func emptySeparatorHandler(content string) ([]string, error) {
	r := strings.NewReader(content)
	buf := make([]byte, maxCharacterPerMsg-1)
	out := []string{}
	for {
		n, err := r.Read(buf)
		if n != 0 {
			out = append(out, string(buf[:n]))
		}

		switch {
		case err == io.EOF:
			return out, nil
		case err != nil:
			return nil, gerrors.Augment(err, "failed_to_batch_emtpy_seperator", nil)
		}
	}
}

func nonEmptySeparatorHandler(content, separator string) ([]string, error) {
	lines := strings.Split(content, separator)

	var (
		total int
		sb    strings.Builder
		out   []string
	)
	for _, line := range lines {
		if len(line) > maxCharacterPerMsg {
			return nil, gerrors.BadParam("bad_param.content.line_too_large", map[string]string{
				"line": line,
			})
		}

		lineLen := len(line) + 1 // +1 to include the newline that the split removed.
		total += lineLen

		switch {
		case total < maxCharacterPerMsg:
			sb.WriteString(fmt.Sprintf("%s\n", line))
		default:
			out = append(out, sb.String())
			sb.Reset()
			total = 0
			sb.WriteString(fmt.Sprintf("%s\n", line))
		}
	}

	// Flush the buffer to capture the last line.
	if sb.String() != "" {
		out = append(out, sb.String())
	}

	return out, nil
}
