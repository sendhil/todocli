package cmd

import (
	"fmt"
	"os"

	"github.com/sendhil/todocli/pkg/todocli"
)

func retrieveItemsWithMetadata() []todocli.Todo {
	retriever := todocli.NewTodoRetriever()
	items, err := retriever.GetItemsWithMetadata()
	if err != nil {
		fmt.Println("Error : ", err)
		os.Exit(1)
	}

	return items
}
