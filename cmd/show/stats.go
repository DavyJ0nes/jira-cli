package show

import (
	"fmt"
	"math"
	"time"

	"github.com/spf13/cobra"

	"github.com/davyj0nes/jira-cli/config"
	"gopkg.in/andygrunwald/go-jira.v1"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "show stats on current sprint",
	Run: func(cmd *cobra.Command, args []string) {
		stats()
	},
}

type issueInfo struct {
	ID                      string
	Summary                 string
	TimeSinceCreated        string
	TimeSinceCreatedSeconds float64
	Assignee                string
}

type issueStats struct {
	Open       []issueInfo
	InProgress []issueInfo
	Review     []issueInfo
	Closed     []issueInfo
	Total      int
}

func stats() {
	j := config.NewJIRAClient().Connect()

	var activeSprint jira.Sprint
	for _, s := range j.GetSprints() {
		if s.State == "active" {
			activeSprint = s
		}
	}
	if activeSprint.ID == 0 {
		panic("no active sprint found")
	}

	issues, _, err := j.Client.Sprint.GetIssuesForSprint(activeSprint.ID)
	if err != nil {
		panic(err)
	}

	stats := parseStats(issues)

	percentComplete := (float64(len(stats.Closed)) / float64(stats.Total)) * 100
	percentExpectedComplete := (float64(len(stats.Closed)+len(stats.Review)) / float64(stats.Total)) * 100
	sprintEnds := time.Until(time.Time(*activeSprint.EndDate)).Round(time.Hour).Seconds()
	sprintEndsDays := math.Floor(sprintEnds / (60 * 60 * 24))
	sprintEndsRemainder := sprintEnds - sprintEndsDays*(60*60*24)
	sprintEndsHour := sprintEndsRemainder / (60 * 60)

	fmt.Println("# Sprint Stats")
	fmt.Println("  ## Stories")
	fmt.Println("     Total:\t\t\t", stats.Total)
	fmt.Println("     Open:\t\t\t", len(stats.Open))
	fmt.Println("     In Progress:\t\t", len(stats.InProgress))
	fmt.Println("     In Review:\t\t\t", len(stats.Review))
	fmt.Println("     Closed:\t\t\t", len(stats.Closed))
	fmt.Println("---------------------------------------")
	fmt.Printf("     %%age Complete:\t\t %.1f%%\n", percentComplete)
	fmt.Printf("     Expected %%age Complete:\t %.1f%%\n", percentExpectedComplete)
	fmt.Println("---------------------------------------")
	fmt.Println("  ## Timing")
	fmt.Println("     Sprint Started:\t\t", activeSprint.StartDate)
	fmt.Println("     Sprint Ends:\t\t", activeSprint.EndDate)
	fmt.Printf("     Remaining Time:\t\t %.0f days %vh\n", sprintEndsDays, sprintEndsHour)
	fmt.Println("---------------------------------------")

}

func parseStats(issues []jira.Issue) issueStats {
	stats := issueStats{}
	for _, issue := range issues {
		stats.Total++
		assigneeName := "Unassigned"
		if issue.Fields.Assignee != nil {
			assigneeName = issue.Fields.Assignee.DisplayName
		}

		createdTime := time.Time(issue.Fields.Created)
		timeSinceCreated := time.Since(createdTime).Round(time.Second)

		info := issueInfo{
			ID:                      issue.Key,
			TimeSinceCreated:        timeSinceCreated.String(),
			TimeSinceCreatedSeconds: timeSinceCreated.Seconds(),
			Summary:                 issue.Fields.Summary,
			Assignee:                assigneeName,
		}

		switch issue.Fields.Status.Name {
		case "In Progress":
			stats.InProgress = append(stats.InProgress, info)
		case "Review":
			stats.Review = append(stats.Review, info)
		case "Resolved":
			stats.Closed = append(stats.Closed, info)
		default:
			stats.Open = append(stats.Open, info)
		}
	}

	return stats
}
