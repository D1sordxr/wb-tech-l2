package command

import "time"

type UpdateEvent struct {
	ID     string    `json:"id"`
	UserID string    `json:"user_id"`
	Date   time.Time `json:"date"`
	Event  string    `json:"event"`
}
