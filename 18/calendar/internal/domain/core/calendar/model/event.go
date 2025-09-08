package model

import (
	"time"

	"github.com/google/uuid"
)

type CalendarEvent struct {
	ID        uuid.UUID `json:"id,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Event     string    `json:"event,omitempty"`
	Date      time.Time `json:"date,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewCalendarEvent(
	ID uuid.UUID,
	userID uuid.UUID,
	event string,
	date time.Time,
) *CalendarEvent {
	return &CalendarEvent{
		ID:        ID,
		UserID:    userID,
		Event:     event,
		Date:      date,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
