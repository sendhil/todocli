package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sendhil/todocli/pkg/todocli"
)

func main() {
	retriever := todocli.NewTodoRetriever()
	items, err := retriever.GetItemsWithMetadata()
	if err != nil {
		panic(err)
	}
	spew.Dump(items)
}
