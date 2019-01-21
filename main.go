package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/structtag"
	"github.com/go-yaml/yaml"
	"github.com/sendhil/todocli/pkg/todocli"
)

func main() {
	fmt.Println(getPathOfTodoItems())
	items, err := parseRawTodoItems(getRawTodoItems())
	if err != nil {
		panic(err)
	}
	spew.Dump(items)
}

func getPathOfTodoItems() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/.todocli.yaml", usr.HomeDir))
	if err != nil {
		panic(err)
	}

	config := todocli.Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config.Path
}

func getRawTodoItems() []string {
	out, err := exec.Command("rg", "-i", "\\[ \\].*`.*`", getPathOfTodoItems()).Output()
	if err != nil {
		panic(err)
	}
	return strings.Split(string(out), "\n")
}

func parseRawTodoItems(rawItems []string) ([]todocli.Todo, error) {
	parsedItems := make([]todocli.Todo, 0)

	todoItemRegex := regexp.MustCompile("\\[ \\]\\s?([^`]*)`(.*)`")

	for _, rawItem := range rawItems {
		if rawItem == "" {
			continue
		}
		// 2. Extract metadata
		result := todoItemRegex.FindAllStringSubmatch(rawItem, -1)
		todoItemText := result[0][1]
		todoItemMetadata := result[0][2]

		item, err := attachMetadata(todocli.Todo{Text: todoItemText}, todoItemMetadata)
		if err != nil {
			return parsedItems, err
		}

		parsedItems = append(parsedItems, item)
	}

	return parsedItems, nil
}

func attachMetadata(item todocli.Todo, metadata string) (todocli.Todo, error) {
	itemToReturn := todocli.Todo{Text: item.Text}

	fmt.Println("Attempting to parse :", metadata)
	tags, err := structtag.Parse(metadata)
	if err != nil {
		return itemToReturn, err
	}

	// 1. Extract Due Date
	dueDateTag, err := tags.Get("due")
	if err == nil {
		// The only error that comes out seems to indicate that the tag doesn't exist
		dueDate, err := time.Parse("01-02-2006", dueDateTag.Name)
		if err != nil {
			return itemToReturn, err
		}
		itemToReturn.Due = dueDate
	}

	// 2. Extract Created At
	createdAt, err := tags.Get("created_at")
	if err == nil {
		// The only error that comes out seems to indicate that the tag doesn't exist
		createdAt, err := time.Parse("01-02-2006", createdAt.Name)
		if err != nil {
			return itemToReturn, err
		}
		itemToReturn.CreatedAt = createdAt
	}

	// 3. Extract Modified At
	modifiedAt, err := tags.Get("modified_at")
	if err == nil {
		// The only error that comes out seems to indicate that the tag doesn't exist
		modifiedAt, err := time.Parse("01-02-2006", modifiedAt.Name)
		if err != nil {
			return itemToReturn, err
		}
		itemToReturn.ModifiedAt = modifiedAt
	}

	// 4. Determine tags and determine if this task has been flagged as important
	tagData, err := tags.Get("tag")
	if err == nil {
		if strings.Contains(strings.ToLower(tagData.Name), "important") {
			itemToReturn.Important = true
		}
	}

	return itemToReturn, nil
}
