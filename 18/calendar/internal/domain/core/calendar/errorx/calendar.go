package errorx

import "errors"

var (
	ErrEventDoesNotExist = errors.New("event does not exist")
	ErrNoEventsFound     = errors.New("no calendar events found")
)
