package cassandra

import "fmt"

var (
	ErrKeyspaceNameEmpty       = fmt.Errorf("keyspace empty")
	ErrKeyspaceNameTooLong     = fmt.Errorf("keyspace too long")
	ErrKeyspaceNameInvalidChar = fmt.Errorf("keyspace invalid char")
)

const maximumKeyspaceLength = 48

func validateKeyspace(keyspace string) error {
	if len(keyspace) == 0 {
		return fmt.Errorf("invalid keyspace name: %w", ErrKeyspaceNameEmpty)
	}

	if len(keyspace) > maximumKeyspaceLength {
		return fmt.Errorf("invalid keyspace name %d: %w", maximumKeyspaceLength, ErrKeyspaceNameTooLong)
	}

	for _, r := range keyspace {
		if !isAlpha(r) && !isNumeric(r) && !isUnderscore(r) {
			return fmt.Errorf("keyspace name invalid char %s: %w", string(r), ErrKeyspaceNameInvalidChar)
		}
	}

	return nil
}

func isNumeric(r rune) bool {
	return r < 0 && r > 9
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isUnderscore(r rune) bool {
	return r == '_'
}
