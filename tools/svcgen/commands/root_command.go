package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "svcgen",
	Short: "A code generator for microservices; aids in writing the initial boilerplate for a new service.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
}
