package show

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/davyj0nes/jira-cli/config"
	"gopkg.in/andygrunwald/go-jira.v1"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "show all JIRA Issues",
	Run: func(cmd *cobra.Command, args []string) {
		showAll()
	},
}

func showAll() {
	j := config.NewJIRAClient().Connect()

	var activeSprint int
	for _, s := range j.GetSprints() {
		if s.State == "active" {
			activeSprint = s.ID
		}
	}
	if activeSprint == 0 {
		panic("no active sprint found")
	}

	issues, _, err := j.Client.Sprint.GetIssuesForSprint(activeSprint)
	if err != nil {
		panic(err)
	}

	printIssues(issues)

}

func printIssues(issues []jira.Issue) {
	for _, issue := range issues {
		assigneeName := "Unassigned"
		if issue.Fields.Assignee != nil {
			assigneeName = issue.Fields.Assignee.DisplayName
		}

		fmt.Println("---------------------")
		fmt.Println("ID:", issue.Key)
		fmt.Println("Summary:", issue.Fields.Summary)
		fmt.Println("Assignee:", assigneeName)
		fmt.Println("Status:", issue.Fields.Status.Name)
	}
}
