package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sendhil/todocli/pkg/todocli"
	"github.com/sendhil/todocli/pkg/todocli/utils"
)

func retrieveItems() []todocli.Todo {
	retriever := todocli.NewTodoRetriever()
	items, err := retriever.GetItems()
	if err != nil {
		fmt.Println("Error : ", err)
		os.Exit(1)
	}

	return items
}

func retrieveItemsWithMetadata() []todocli.Todo {
	retriever := todocli.NewTodoRetriever()
	items, err := retriever.GetItemsWithMetadata()
	if err != nil {
		fmt.Println("Error : ", err)
		os.Exit(1)
	}

	return items
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
