package main

import (
	"fmt"
	"os"

	"github.com/davyj0nes/jira-cli/cmd/show"

	"github.com/davyj0nes/jira-cli/config"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jli",
	Short: "Simple CLI for working with JIRA",
	Long: `simple CLI tool for interaction with JIRA
	on a daily basis.`,
}

func init() {
	rootCmd.AddCommand(show.ShowCmd)
}

func main() {
	config.LoadConfigFile()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
