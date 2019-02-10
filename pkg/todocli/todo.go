package todocli

import "time"

// Todo is the structure that holds all the data pertaining to a Todo item
type Todo struct {
	Text       string
	Filename   string
	Tag        string
	Important  bool
	Due        time.Time
	CreatedAt  time.Time
	ModifiedAt time.Time
}
