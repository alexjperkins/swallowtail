package util

import (
	"fmt"
	"os"
)

// SetEnv wrapper around `os.Lookup` to provide safety if env var is missing.
func SetEnv(variableName string) string {
	v, ok := os.LookupEnv(variableName)
	if !ok {
		panic(fmt.Sprintf("%s: not found in env", variableName))
	}
	return v
}
