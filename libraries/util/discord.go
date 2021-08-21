package util

import "fmt"

// WrapAsCodeBlock wraps the content as a code block.
func WrapAsCodeBlock(content string) string {
	return fmt.Sprintf("```%s```", content)
}
