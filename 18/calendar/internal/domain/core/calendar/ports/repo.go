package ports

import (
	"context"
	"time"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/model"

	"github.com/google/uuid"
)

type CalendarEventRepo interface {
	basicRepo
	filterRepo
}

type basicRepo interface {
	Save(ctx context.Context, event *model.CalendarEvent) error
	Read(ctx context.Context, userID, eventID uuid.UUID) (*model.CalendarEvent, error)
	Update(ctx context.Context, event *model.CalendarEvent) error
	Delete(ctx context.Context, userID, eventID uuid.UUID) error
}

type filterRepo interface {
	ReadByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*model.CalendarEvent, error)
	ReadBetweenDates(ctx context.Context, userID uuid.UUID, dateStart, dateEnd time.Time) ([]*model.CalendarEvent, error)
}
