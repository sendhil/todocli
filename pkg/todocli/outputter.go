package todocli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/sendhil/todocli/pkg/todocli/utils"
)

// Outputter interface allows the outputting of Todo items in a nicer to look format
type Outputter interface {
	OutputTodoItems(items []Todo)
}

type outputter struct {
}

func (o *outputter) OutputTodoItems(items []Todo) {
	today := time.Now()
	year, month, day := today.Date()
	weekday := today.Weekday()
	noDateTime := time.Time{}.Add(time.Minute)

	itemsMap := make(map[int]Todo)
	for index, item := range items {
		itemsMap[index] = item
	}

	itemsByBucket := make(map[string][]Todo, 0)
	itemsByBucket["No Date"] = getItemsBeforeDate(itemsMap, noDateTime)

	startOfToday := time.Date(year, month, day, 0, 0, 0, 0, today.Location()).Add(time.Minute * -1)
	itemsByBucket["The Past"] = getItemsBeforeDate(itemsMap, startOfToday)

	itemsByBucket["Today"] = getItemsBeforeDate(itemsMap, time.Date(year, month, day, 0, 0, 0, 0, today.Location()).Add(time.Hour*24).Add(time.Minute*-1))

	var daysToAdjust int
	if weekday == time.Sunday {
		daysToAdjust = 7
	} else {
		daysToAdjust = (int)(time.Saturday-weekday) + 1
	}
	endOfWeek := time.Date(year, month, day, 0, 0, 0, 0, today.Location()).Add(time.Hour * 24 * time.Duration(daysToAdjust)).Add(-1 * time.Minute)
	itemsByBucket["This Week"] = getItemsBeforeDate(itemsMap, endOfWeek)
	itemsByBucket["Farther Out"] = getItemsBeforeDate(itemsMap, time.Date(3000, 01, 01, 0, 0, 0, 0, today.Location()))

	itemsPrinted := 0
	printedSpacer := false

	config := utils.GetConfig()

	for key, items := range itemsByBucket {
		if len(items) == 0 {
			continue
		}

		if !printedSpacer {
			printedSpacer = true
			fmt.Println("")
		}

		fmt.Printf("%s:\n", key)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Index", "File", "Description", "Due", "Tag", "Important"})

		for index, item := range items {
			indexAsString := strconv.Itoa(index + 1)
			importantAsString := strconv.FormatBool(item.Important)
			table.Append([]string{indexAsString, getFriendlyFilename(item.Filename, config), item.Text, item.Due.Format("01/02/2006"), item.Tag, importantAsString})
		}

		table.Render()
		fmt.Println("")

		itemsPrinted += len(items)
	}

	if len(items) != itemsPrinted {
		fmt.Printf("Expected to print : %d but instead printed %d", len(items), itemsPrinted)
	}
}

func getItemsBeforeDate(itemsMap map[int]Todo, endDate time.Time) []Todo {
	itemsToReturn := make([]Todo, 0)
	itemsToRemove := make([]int, 0)
	for index, item := range itemsMap {
		if endDate.After(item.Due) {
			itemsToRemove = append(itemsToRemove, index)
			itemsToReturn = append(itemsToReturn, item)
		}
	}

	for _, itemIndex := range itemsToRemove {
		delete(itemsMap, itemIndex)
	}

	return itemsToReturn
}

func getFriendlyFilename(fileName string, config *utils.Config) string {
	if val, ok := config.Mappings[fileName]; ok {
		return val
	}

	lastIndexOfSlash := strings.LastIndex(fileName, "/")
	if lastIndexOfSlash > 0 {
		return fileName[lastIndexOfSlash+1 : len(fileName)]
	}

	return fileName
}

// NewOutputter creates an object to output Todo items
func NewOutputter() Outputter {
	return &outputter{}
}
