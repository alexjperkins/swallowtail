package util

import (
	"os"
)

// SetEnv wrapper around `os.Lookup` to provide safety if env var is missing.
func SetEnv(variableName string) string {
	v, _ := os.LookupEnv(variableName)
	return v
}

func EnvGetOrDefault(key string, d interface{}) interface{} {
	v, ok := os.LookupEnv(key)
	if !ok {
		return d
	}
	return v
}
