package command

import "time"

type CreateEvent struct {
	UserID string    `json:"user_id"`
	Event  string    `json:"event"`
	Date   time.Time `json:"date"`
}
