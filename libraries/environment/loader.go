package environment

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadEnvironment loads the environment and returns as a typed struct.
//
// Returns an error if for some reason we fail to read the environment file.
func LoadEnvironment() (*Environment, error) {
	// Validate environent file env var.
	envFile := os.Getenv(environmentFileName)
	if envFile == "" {
		return nil, fmt.Errorf("load environment variable: %w", ErrMissingEnvironmentFileEnvVar)
	}

	// Load environment.
	if err := loadEnvFile(envFile); err != nil {
		return nil, fmt.Errorf("load environment: %w", err)
	}

	// Process environment.
	var env *Environment
	if err := envconfig.Process("", env); err != nil {
		return nil, fmt.Errorf("process environment: %w", err)
	}

	return env, nil
}

func loadEnvFile(filename string) error {
	// Validate file.
	if _, err := os.Stat(filename); err != nil {
		return fmt.Errorf("load envfile; missing: %w", err)
	}

	// Load environment file.
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("load envfile: %w", err)
	}

	return nil
}
