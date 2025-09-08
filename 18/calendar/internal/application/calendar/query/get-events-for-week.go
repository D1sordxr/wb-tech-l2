package query

import "time"

type GetEventsForWeek struct {
	UserID    string    `json:"user_id"`
	DateStart time.Time `json:"date_start"`
}
