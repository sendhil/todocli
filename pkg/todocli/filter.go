package todocli

import "time"

// TODO: Revisit name of this interface

// Filter filters todo items by field
type Filter interface {
	GetItemsBetweenDates(items []Todo, startDate, endDate time.Time) []Todo
	GetImportantItems(items []Todo) []Todo
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

	return itemsToReturn
}

func (f *filter) GetImportantItems(items []Todo) []Todo {
	itemsToReturn := make([]Todo, 0)

	for _, item := range items {
		if item.Important {
			itemsToReturn = append(itemsToReturn, item)
		}
	}

	return itemsToReturn
}

// NewFilter returns an object that can Filter Todo items
func NewFilter() Filter {
	return &filter{}
}
