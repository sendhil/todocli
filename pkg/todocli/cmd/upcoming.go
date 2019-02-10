package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sendhil/todocli/pkg/todocli"
	"github.com/sendhil/todocli/pkg/todocli/utils"
	"github.com/spf13/cobra"
)

var upcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "Print out tasks that are upcoming",
	Run: func(cmd *cobra.Command, args []string) {
		upcomingAllCmd.Run(cmd, args)
	},
}

var upcomingAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Print out all tasks that are upcoming",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		year, month, day := t.Date()
		startOfToday := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
		dateFartherOut := time.Date(3000, 01, 01, 0, 0, 0, 0, t.Location())
		retrieveAndPrintTasksByDate(startOfToday, dateFartherOut)
	},
}

var upcomingTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "Print out tasks that are upcoming today",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		year, month, day := t.Date()
		startOfToday := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
		endOfToday := time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Add(24 * time.Hour).Add(-1 * time.Minute)
		retrieveAndPrintTasksByDate(startOfToday, endOfToday)
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
		retrieveAndPrintTasksByDate(startOfWeek, endOfWeek)
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
		retrieveAndPrintTasksByDate(startOfMonth, endOfMonth)
	},
}

func retrieveAndPrintTasksByDate(startDate, endDate time.Time) {
	if Verbose {
		fmt.Printf("Retrieving items between %v and %v\n", startDate, endDate)
	}

	items := retrieveItemsWithMetadata()

	filter := todocli.NewFilter()
	filteredItems := filter.GetItemsBetweenDates(items, startDate, endDate)
	filteredItems = filter.GetItemsWithTag(filteredItems, Tag)
	filteredItems = filterItemsByFile(filter, filteredItems)

	outputter := todocli.NewOutputter()
	outputter.OutputTodoItems(filteredItems)
}

// Attempts to filter by either the filename passed in or the file alias.
func filterItemsByFile(filter todocli.Filter, items []todocli.Todo) []todocli.Todo {
	fileName := File
	matched := false
	if len(FileAlias) > 0 {
		config := utils.GetConfig()
		for key, value := range config.Mappings {
			if strings.ToLower(FileAlias) == strings.ToLower(value) {
				fileName = key
				matched = true
				break
			}
		}
	}

	if len(FileAlias) > 0 && !matched {
		color.Yellow(fmt.Sprintf("Warning: Unable to match alias '%s'", FileAlias))
		return []todocli.Todo{}
	}

	return filter.GetItemsWithFile(items, fileName)
}

func init() {
	upcomingCmd.AddCommand(upcomingAllCmd)
	upcomingCmd.AddCommand(upcomingTodayCmd)
	upcomingCmd.AddCommand(upcomingWeekCmd)
	upcomingCmd.AddCommand(upcomingMonthCmd)
	rootCmd.AddCommand(upcomingCmd)
}
