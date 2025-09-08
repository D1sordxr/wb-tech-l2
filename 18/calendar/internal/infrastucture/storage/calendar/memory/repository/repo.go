package repository

import (
	"context"
	"fmt"
	"time"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/errorx"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/model"
	"wb-tech-l2/18/calendar/internal/infrastucture/storage/calendar/memory"

	"github.com/google/uuid"
)

type CalendarRepo struct {
	storage *memory.CalendarStorage
}

func NewCalendarRepo(storage *memory.CalendarStorage) *CalendarRepo {
	return &CalendarRepo{storage: storage}
}

func (c *CalendarRepo) Save(ctx context.Context, event *model.CalendarEvent) error {
	return c.storage.Set(ctx, event.UserID.String(), event)
}

func (c *CalendarRepo) Read(ctx context.Context, userID, eventID uuid.UUID) (*model.CalendarEvent, error) {
	return c.storage.Get(ctx, userID.String(), eventID.String())
}

func (c *CalendarRepo) Update(ctx context.Context, event *model.CalendarEvent) error {
	if c.storage.Exists(ctx, event.UserID.String(), event.ID.String()) {
		return c.storage.Set(ctx, event.UserID.String(), event)
	}

	return errorx.ErrEventDoesNotExist
}

func (c *CalendarRepo) Delete(ctx context.Context, userID, eventID uuid.UUID) error {
	return c.storage.Delete(ctx, userID.String(), eventID.String())
}

func (c *CalendarRepo) ReadByDate(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
) ([]*model.CalendarEvent, error) {
	const op = "storage.memory.CalendarRepo.ReadByDate"

	allEvents, err := c.storage.List(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	events := make([]*model.CalendarEvent, 0, 4)
	endOfDay := date.Add(24 * time.Hour)
	for _, event := range allEvents {
		if (event.Date.Equal(date) || event.Date.After(date)) &&
			event.Date.Before(endOfDay) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (c *CalendarRepo) ReadBetweenDates(
	ctx context.Context,
	userID uuid.UUID,
	dateStart, dateEnd time.Time,
) ([]*model.CalendarEvent, error) {
	const op = "storage.memory.CalendarRepo.ReadBetweenDates"

	allEvents, err := c.storage.List(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	events := make([]*model.CalendarEvent, 0, 16)
	for _, event := range allEvents {
		if (event.Date.Equal(dateStart) || event.Date.After(dateStart)) &&
			event.Date.Before(dateEnd) {
			events = append(events, event)
		}
	}

	return events, nil
}
