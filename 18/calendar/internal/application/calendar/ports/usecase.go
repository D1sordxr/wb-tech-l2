package ports

import (
	"context"
	"wb-tech-l2/18/calendar/internal/application/calendar/command"
	"wb-tech-l2/18/calendar/internal/application/calendar/query"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/model"
)

type CalendarUseCase interface {
	commands
	queries
}

type commands interface {
	CreateEvent(ctx context.Context, cmd command.CreateEvent) (string, error)
	UpdateEvent(ctx context.Context, cmd command.UpdateEvent) error
	DeleteEvent(ctx context.Context, cmd command.DeleteEvent) error
}

type queries interface {
	GetEventsForDay(ctx context.Context, q query.GetEventsForDay) ([]*model.CalendarEvent, error)
	GetEventsForWeek(ctx context.Context, q query.GetEventsForWeek) ([]*model.CalendarEvent, error)
	GetEventsForMonth(ctx context.Context, q query.GetEventsForMonth) ([]*model.CalendarEvent, error)
}
