package cmd

import (
	"github.com/sendhil/todocli/pkg/todocli"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Print out all tasks with metadata",
	Run: func(cmd *cobra.Command, args []string) {
		retrieveAllTasks()
		upcomingAllCmd.Run(cmd, args)
	},
}

func retrieveAllTasks() {
	items := retrieveItems()

	filter := todocli.NewFilter()
	items = filter.GetItemsWithTag(items, Tag)

	outputter := todocli.NewOutputter()
	outputter.OutputTodoItems(items)
}

func init() {
	rootCmd.AddCommand(allCmd)
}
