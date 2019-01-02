package show

import (
	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show JIRA Issues",
}

func init() {
	ShowCmd.AddCommand(allCmd)
	ShowCmd.AddCommand(statsCmd)
}
