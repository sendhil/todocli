package todocli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/structtag"
	"github.com/go-yaml/yaml"
)

// TodoRetriever represents an interface that can retrieve Todo Items
type TodoRetriever interface {
	// GetItems retrieves all opened []Todo items
	GetItems() ([]Todo, error)

	// GetItemsWithMetadata retrieves all opened []Todo items that have metadata associated
	GetItemsWithMetadata() ([]Todo, error)
}

type todoRetriever struct {
}

func (t *todoRetriever) GetItems() ([]Todo, error) {
	return parseRawTodoItems(getRawTodoItems())
}

func (t *todoRetriever) GetItemsWithMetadata() ([]Todo, error) {
	return parseRawTodoItems(getRawTodoItemsWithMetadata())
}

// NewTodoRetriever creates an object that can retrieve Todo item data
func NewTodoRetriever() TodoRetriever {
	return &todoRetriever{}
}

func getRawTodoItems() []string {
	out, err := exec.Command("rg", "-i", "\\[ \\]", getPathOfTodoItems()).Output()
	if err != nil {
		panic(err)
	}
	return strings.Split(string(out), "\n")
}

func getRawTodoItemsWithMetadata() []string {
	out, err := exec.Command("rg", "-i", "\\[ \\].*`.*`", getPathOfTodoItems()).Output()
	if err != nil {
		panic(err)
	}
	return strings.Split(string(out), "\n")
}

func getPathOfTodoItems() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configPath := fmt.Sprintf("%s/.todocli.yaml", usr.HomeDir)

	if _, err := os.Stat(configPath); err != nil {
		panic("Can't find .todocli.yaml")
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config.Path
}

func parseRawTodoItems(rawItems []string) ([]Todo, error) {
	parsedItems := make([]Todo, 0)

	todoItemRegex := regexp.MustCompile("\\[ \\]\\s?([^`]*)`(.*)`")

	for _, rawItem := range rawItems {
		if rawItem == "" {
			continue
		}
		// 2. Extract metadata
		result := todoItemRegex.FindAllStringSubmatch(rawItem, -1)

		// If there's no metadata then just attach the raw text
		if len(result) == 0 {
			parsedItems = append(parsedItems, Todo{Text: rawItem})
			continue
		}

		todoItemText := result[0][1]
		todoItemMetadata := result[0][2]

		item, err := attachMetadata(Todo{Text: todoItemText}, todoItemMetadata)
		if err != nil {
			return parsedItems, err
		}

		parsedItems = append(parsedItems, item)
	}

	return parsedItems, nil
}

func attachMetadata(item Todo, metadata string) (Todo, error) {
	itemToReturn := Todo{Text: item.Text}

	tags, err := structtag.Parse(metadata)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error parsing Metdata : '%s' for '%s'", metadata, item.Text))
		return itemToReturn, err
	}

	// 1. Extract Due Date
	dueDateTag, err := tags.Get("due")
	if err == nil {
		// The only error that comes out seems to indicate that the tag doesn't exist
		dueDate, err := time.ParseInLocation("01-02-2006", dueDateTag.Name, time.Now().Location())
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

		if len(tagData.Name) > 0 {
			itemToReturn.Tag = tagData.Name
		}
	}

	return itemToReturn, nil
}
