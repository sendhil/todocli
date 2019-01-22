package cmd

import (
	"github.com/sendhil/todocli/pkg/todocli"
	"github.com/spf13/cobra"
)

var importantCmd = &cobra.Command{
	Use:   "important",
	Short: "Print out tasks that are important",
	Run: func(cmd *cobra.Command, args []string) {
		retrieveAndPrintImportantTasks()
	},
}

func init() {
	rootCmd.AddCommand(importantCmd)
}

func retrieveAndPrintImportantTasks() {
	items := retrieveItemsWithMetadata()

	filter := todocli.NewFilter()
	filteredItems := filter.GetImportantItems(items)

	outputter := todocli.NewOutputter()
	outputter.OutputTodoItems(filteredItems)
}
