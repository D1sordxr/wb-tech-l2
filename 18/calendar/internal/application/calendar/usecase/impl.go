package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
	"wb-tech-l2/18/calendar/internal/application/calendar/command"
	"wb-tech-l2/18/calendar/internal/application/calendar/query"
	appPorts "wb-tech-l2/18/calendar/internal/domain/app/ports"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/errorx"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/model"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/ports"
)

type Calendar struct {
	log  appPorts.Logger
	repo ports.CalendarEventRepo
}

func NewCalendar(
	log appPorts.Logger,
	repo ports.CalendarEventRepo,
) *Calendar {
	return &Calendar{
		log:  log,
		repo: repo,
	}
}

func (c *Calendar) CreateEvent(ctx context.Context, cmd command.CreateEvent) (string, error) {
	const op = "application.calendar.UseCase.CreateEvent"
	withFields := func(fields ...any) []any {
		return append([]any{"op", op, "user_id", cmd.UserID}, fields...)
	}

	c.log.Info("Attempting to create calendar event", withFields()...)

	eventID := uuid.New()
	userID, err := c.parseUserID(cmd.UserID, op)
	if err != nil {
		return "", err
	}

	if err = c.repo.Save(ctx, model.NewCalendarEvent(
		eventID,
		userID,
		cmd.Event,
		cmd.Date,
	)); err != nil {
		c.log.Error("Failed to save calendar event", withFields("error", err.Error())...)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Successfully created calendar event", withFields()...)

	return eventID.String(), nil
}

func (c *Calendar) UpdateEvent(ctx context.Context, cmd command.UpdateEvent) error {
	const op = "application.calendar.UseCase.UpdateEvent"
	withFields := func(fields ...any) []any {
		return append([]any{"op", op, "user_id", cmd.UserID, "event_id", cmd.ID}, fields...)
	}

	c.log.Info("Attempting to update calendar event", withFields()...)

	eventID, err := c.parseEventID(cmd.ID, op)
	if err != nil {
		return err
	}
	userID, err := c.parseUserID(cmd.UserID, op)
	if err != nil {
		return err
	}

	if err = c.repo.Update(ctx, model.NewCalendarEvent(
		eventID,
		userID,
		cmd.Event,
		cmd.Date,
	)); err != nil {
		c.log.Error("Failed to save calendar event", withFields("error", err.Error())...)
		return fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Successfully updated calendar event", withFields()...)

	return nil
}

func (c *Calendar) DeleteEvent(ctx context.Context, cmd command.DeleteEvent) error {
	const op = "application.calendar.UseCase.DeleteEvent"
	withFields := func(fields ...any) []any {
		return append([]any{"op", op, "user_id", cmd.UserID, "event_id", cmd.ID}, fields...)
	}

	c.log.Info("Attempting to delete calendar event", withFields()...)

	eventID, err := c.parseEventID(cmd.ID, op)
	if err != nil {
		return err
	}
	userID, err := c.parseUserID(cmd.UserID, op)
	if err != nil {
		return err
	}

	if err = c.repo.Delete(ctx, userID, eventID); err != nil {
		c.log.Error("Failed to delete calendar event", withFields("error", err.Error())...)
		return fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Successfully deleted calendar event", withFields()...)

	return nil
}

func (c *Calendar) GetEventsForDay(ctx context.Context, q query.GetEventsForDay) ([]*model.CalendarEvent, error) {
	const op = "application.calendar.UseCase.GetEventsForDay"
	withFields := setupLogFieldsFunc("op", op, "user_id", q.UserID, "date", q.Date.String())

	c.log.Info("Attempting to get calendar events for day", withFields()...)

	userID, err := c.parseUserID(q.UserID, op)
	if err != nil {
		return nil, err
	}

	events, err := c.repo.ReadByDate(ctx, userID, q.Date)
	if err != nil {
		c.log.Error("Failed to get calendar events", withFields("error", err.Error())...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(events) < 1 {
		c.log.Info("No calendar events found", withFields()...)
		return events, fmt.Errorf("%s: %w", op, errorx.ErrNoEventsFound)
	}

	c.log.Info("Successfully got calendar events for day", withFields()...)

	return events, nil
}

func (c *Calendar) GetEventsForWeek(ctx context.Context, q query.GetEventsForWeek) ([]*model.CalendarEvent, error) {
	const op = "application.calendar.UseCase.GetEventsForWeek"
	withFields := setupLogFieldsFunc("op", op, "user_id", q.UserID, "date_start", q.DateStart.String())

	c.log.Info("Attempting to get calendar events for week", withFields()...)

	userID, err := c.parseUserID(q.UserID, op)
	if err != nil {
		return nil, err
	}

	dateEnd := q.DateStart.Add(time.Hour * 24 * 7)

	events, err := c.repo.ReadBetweenDates(ctx, userID, q.DateStart, dateEnd)
	if err != nil {
		c.log.Error("Failed to get calendar events", withFields("error", err.Error())...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(events) < 1 {
		c.log.Info("No calendar events found", withFields()...)
		return events, fmt.Errorf("%s: %w", op, errorx.ErrNoEventsFound)
	}

	c.log.Info("Successfully got calendar events for week", withFields()...)

	return events, nil
}

func (c *Calendar) GetEventsForMonth(ctx context.Context, q query.GetEventsForMonth) ([]*model.CalendarEvent, error) {
	const op = "application.calendar.UseCase.GetEventsForMonth"
	withFields := setupLogFieldsFunc("op", op, "user_id", q.UserID, "date_start", q.DateStart.String())

	c.log.Info("Attempting to get calendar events for month", withFields()...)

	userID, err := c.parseUserID(q.UserID, op)
	if err != nil {
		return nil, err
	}

	dateEnd := q.DateStart.AddDate(0, 1, 0) // q.DateStart.Add(time.Hour * 24 * 30)

	events, err := c.repo.ReadBetweenDates(ctx, userID, q.DateStart, dateEnd)
	if err != nil {
		c.log.Error("Failed to get calendar events", withFields("error", err.Error())...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(events) < 1 {
		c.log.Info("No calendar events found", withFields()...)
		return events, fmt.Errorf("%s: %w", op, errorx.ErrNoEventsFound)
	}

	c.log.Info("Successfully got calendar events for month", withFields()...)

	return events, nil
}

func setupLogFieldsFunc(fields ...any) func(...any) []any {
	return func(newFields ...any) []any {
		return append(fields, newFields...)
	}
}

func (c *Calendar) parseUserID(userID, op string) (uuid.UUID, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		c.log.Error("Failed to parse user ID", "op", op, "user_id", userID, "error", err.Error())
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (c *Calendar) parseEventID(eventID, op string) (uuid.UUID, error) {
	id, err := uuid.Parse(eventID)
	if err != nil {
		c.log.Error("Failed to parse event ID", "op", op, "event_id", eventID, "error", err.Error())
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
