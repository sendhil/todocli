package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/sendhil/todocli/pkg/todocli"
	"github.com/spf13/cobra"
)

var upcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "Print out tasks that are upcoming",
}

var upcomingTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "Print out tasks that are upcoming today",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		year, month, day := t.Date()
		startOfToday := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
		endOfToday := time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Add(24 * time.Hour).Add(-1 * time.Minute)
		retrieveAndPrintTasks(startOfToday, endOfToday)
	},
}

var upcomingWeekCmd = &cobra.Command{
	Use:   "week",
	Short: "Print out tasks that are upcoming this week",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		year, month, day := t.Date()
		weekday := t.Weekday()
		var daysToAdjust int
		if weekday == time.Sunday {
			daysToAdjust = 7
		} else {
			daysToAdjust = (int)(time.Saturday-weekday) + 1
		}
		startOfWeek := time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Add(time.Hour * 24 * time.Duration(daysToAdjust-7))
		endOfWeek := time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Add(time.Hour * 24 * time.Duration(daysToAdjust)).Add(-1 * time.Minute)
		retrieveAndPrintTasks(startOfWeek, endOfWeek)
	},
}

var upcomingMonthCmd = &cobra.Command{
	Use:   "month",
	Short: "Print out tasks that are upcoming this month",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		year, month, _ := t.Date()
		if month == 12 {
			year = year + 1
			month = 0
		}
		startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
		endOfMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location()).Add(-1 * time.Minute)
		retrieveAndPrintTasks(startOfMonth, endOfMonth)
	},
}

func retrieveAndPrintTasks(startDate, endDate time.Time) {
	if Verbose {
		fmt.Printf("Retrieving items between %v and %v\n", startDate, endDate)
	}

	retriever := todocli.NewTodoRetriever()
	items, err := retriever.GetItemsWithMetadata()
	if err != nil {
		fmt.Println("Error : ", err)
		os.Exit(1)
	}

	filter := todocli.NewFilter()
	filteredItems := filter.GetItemsBetweenDates(items, startDate, endDate)

	outputter := todocli.NewOutputter()
	outputter.OutputTodoItems(filteredItems)
}

func init() {
	upcomingCmd.AddCommand(upcomingTodayCmd)
	upcomingCmd.AddCommand(upcomingWeekCmd)
	upcomingCmd.AddCommand(upcomingMonthCmd)
	rootCmd.AddCommand(upcomingCmd)
}
