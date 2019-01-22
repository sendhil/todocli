package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var overdueCmd = &cobra.Command{
	Use:   "overdue",
	Short: "Print out tasks that are overdue",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		year, month, day := t.Date()
		endOfToday := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
		retrieveAndPrintTasksByDate(time.Time{}, endOfToday)
	},
}

func init() {
	rootCmd.AddCommand(overdueCmd)
}
