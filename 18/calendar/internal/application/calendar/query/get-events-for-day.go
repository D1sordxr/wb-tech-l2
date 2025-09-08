package query

import "time"

type GetEventsForDay struct {
	UserID string    `json:"user_id"`
	Date   time.Time `json:"date"`
}
