package todocli

import (
	"sort"
	"strings"
	"time"
)

// TODO: Revisit name of this interface

// Filter filters todo items by field
type Filter interface {
	GetItemsBetweenDates(items []Todo, startDate, endDate time.Time) []Todo
	GetImportantItems(items []Todo) []Todo
	GetItemsWithTag(items []Todo, tag string) []Todo
	GetItemsWithFile(items []Todo, file string) []Todo
}

type filter struct {
}

func (f *filter) GetItemsBetweenDates(items []Todo, startDate, endDate time.Time) []Todo {
	itemsToReturn := make([]Todo, 0)

	for _, item := range items {
		if item.Due.After(startDate) && item.Due.Before(endDate) {
			itemsToReturn = append(itemsToReturn, item)
		}
	}

	sort.Slice(itemsToReturn, func(i, j int) bool {
		if itemsToReturn[i].Due.Equal(itemsToReturn[j].Due) {
			return strings.Compare(itemsToReturn[i].Text, itemsToReturn[j].Text) == -1

		}

		return itemsToReturn[i].Due.Before(itemsToReturn[j].Due)
	})

	return itemsToReturn
}

func (f *filter) GetItemsWithTag(items []Todo, tag string) []Todo {
	if len(tag) == 0 {
		return items
	}

	itemsToReturn := make([]Todo, 0)
	lowerCasedTag := strings.ToLower(tag)

	for _, item := range items {
		for _, tag := range item.Tags {
			if strings.ToLower(tag) == lowerCasedTag {
				itemsToReturn = append(itemsToReturn, item)
			}
		}
	}

	return itemsToReturn
}

func (f *filter) GetItemsWithFile(items []Todo, file string) []Todo {
	if len(file) == 0 {
		return items
	}

	itemsToReturn := make([]Todo, 0)
	lowerCasedFilename := strings.ToLower(file)

	for _, item := range items {
		if strings.ToLower(item.Filename) == lowerCasedFilename {
			itemsToReturn = append(itemsToReturn, item)
		}
	}

	return itemsToReturn
}

func (f *filter) GetImportantItems(items []Todo) []Todo {
	itemsToReturn := make([]Todo, 0)

	for _, item := range items {
		if item.Important {
			itemsToReturn = append(itemsToReturn, item)
		}
	}

	sort.Slice(itemsToReturn, func(i, j int) bool {
		if itemsToReturn[i].Due.Equal(itemsToReturn[j].Due) {
			return strings.Compare(itemsToReturn[i].Text, itemsToReturn[j].Text) == -1

		}

		return itemsToReturn[i].Due.Before(itemsToReturn[j].Due)
	})

	return itemsToReturn
}

// NewFilter returns an object that can Filter Todo items
func NewFilter() Filter {
	return &filter{}
}
