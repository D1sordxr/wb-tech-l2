package usecase

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
	"wb-tech-l2/18/calendar/internal/infrastucture/storage/calendar/memory/repository"

	commands "wb-tech-l2/18/calendar/internal/application/calendar/command"
	queries "wb-tech-l2/18/calendar/internal/application/calendar/query"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/errorx"
	"wb-tech-l2/18/calendar/internal/infrastucture/storage/calendar/memory"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// integration tests

func setupTest() (*Calendar, context.Context) {
	ctx := context.Background()
	storage := memory.NewCalendarStorage()
	repo := repository.NewCalendarRepo(storage)
	logger := slog.Default()

	return NewCalendar(logger, repo), ctx
}

func TestCalendar_CreateEvent(t *testing.T) {
	uc, ctx := setupTest()

	t.Run("successful creation", func(t *testing.T) {
		cmd := commands.CreateEvent{
			UserID: uuid.New().String(),
			Event:  "Test Event",
			Date:   time.Now(),
		}

		eventID, err := uc.CreateEvent(ctx, cmd)

		assert.NoError(t, err)
		assert.NotEmpty(t, eventID)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		cmd := commands.CreateEvent{
			UserID: "invalid-uuid",
			Event:  "Test Event",
			Date:   time.Now(),
		}

		eventID, err := uc.CreateEvent(ctx, cmd)

		assert.Error(t, err)
		assert.Empty(t, eventID)
		assert.Contains(t, err.Error(), "invalid UUID")
	})
}

func TestCalendar_UpdateEvent(t *testing.T) {
	uc, ctx := setupTest()

	t.Run("successful update", func(t *testing.T) {
		userID := uuid.New()
		createCmd := commands.CreateEvent{
			UserID: userID.String(),
			Event:  "Original Event",
			Date:   time.Now(),
		}

		eventID, err := uc.CreateEvent(ctx, createCmd)
		assert.NoError(t, err)

		updateCmd := commands.UpdateEvent{
			ID:     eventID,
			UserID: userID.String(),
			Event:  "Updated Event",
			Date:   time.Now().Add(24 * time.Hour),
		}

		err = uc.UpdateEvent(ctx, updateCmd)
		assert.NoError(t, err)
	})

	t.Run("update non-existent event", func(t *testing.T) {
		cmd := commands.UpdateEvent{
			ID:     uuid.New().String(),
			UserID: uuid.New().String(),
			Event:  "Non-existent Event",
			Date:   time.Now(),
		}

		err := uc.UpdateEvent(ctx, cmd)
		assert.Error(t, err)
	})

	t.Run("invalid event ID", func(t *testing.T) {
		cmd := commands.UpdateEvent{
			ID:     "invalid-uuid",
			UserID: uuid.New().String(),
			Event:  "Test Event",
			Date:   time.Now(),
		}

		err := uc.UpdateEvent(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UUID")
	})
}

func TestCalendar_DeleteEvent(t *testing.T) {
	uc, ctx := setupTest()

	t.Run("successful deletion", func(t *testing.T) {
		userID := uuid.New()
		createCmd := commands.CreateEvent{
			UserID: userID.String(),
			Event:  "Event to delete",
			Date:   time.Now(),
		}

		eventID, err := uc.CreateEvent(ctx, createCmd)
		assert.NoError(t, err)

		deleteCmd := commands.DeleteEvent{
			ID:     eventID,
			UserID: userID.String(),
		}

		err = uc.DeleteEvent(ctx, deleteCmd)
		assert.NoError(t, err)
	})

	t.Run("delete non-existent event", func(t *testing.T) {
		cmd := commands.DeleteEvent{
			ID:     uuid.New().String(),
			UserID: uuid.New().String(),
		}

		err := uc.DeleteEvent(ctx, cmd)
		assert.Error(t, err)
	})
}

func TestCalendar_GetEventsForDay(t *testing.T) {
	uc, ctx := setupTest()

	t.Run("successful get events for day", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now().Truncate(24 * time.Hour)

		createCmd1 := commands.CreateEvent{
			UserID: userID.String(),
			Event:  "Morning Event",
			Date:   now.Add(10 * time.Hour),
		}

		createCmd2 := commands.CreateEvent{
			UserID: userID.String(),
			Event:  "Evening Event",
			Date:   now.Add(18 * time.Hour),
		}

		_, err := uc.CreateEvent(ctx, createCmd1)
		assert.NoError(t, err)
		_, err = uc.CreateEvent(ctx, createCmd2)
		assert.NoError(t, err)

		query := queries.GetEventsForDay{
			UserID: userID.String(),
			Date:   now,
		}

		events, err := uc.GetEventsForDay(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, events, 2)
	})

	t.Run("no events found", func(t *testing.T) {
		query := queries.GetEventsForDay{
			UserID: uuid.New().String(),
			Date:   time.Now(),
		}

		events, err := uc.GetEventsForDay(ctx, query)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errorx.ErrNoEventsFound))
		assert.Empty(t, events)
	})

	t.Run("events from different days", func(t *testing.T) {
		userID := uuid.New()
		today := time.Now().Truncate(24 * time.Hour)
		yesterday := today.Add(-24 * time.Hour)

		createCmd := commands.CreateEvent{
			UserID: userID.String(),
			Event:  "Yesterday Event",
			Date:   yesterday.Add(12 * time.Hour),
		}

		_, err := uc.CreateEvent(ctx, createCmd)
		assert.NoError(t, err)

		query := queries.GetEventsForDay{
			UserID: userID.String(),
			Date:   today,
		}

		events, err := uc.GetEventsForDay(ctx, query)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errorx.ErrNoEventsFound))
		assert.Empty(t, events)
	})
}

func TestCalendar_GetEventsForWeek(t *testing.T) {
	uc, ctx := setupTest()

	t.Run("successful get events for week", func(t *testing.T) {
		userID := uuid.New()
		startOfWeek := time.Now().Truncate(24 * time.Hour)

		for i := 0; i < 3; i++ {
			createCmd := commands.CreateEvent{
				UserID: userID.String(),
				Event:  "Event " + string(rune('A'+i)),
				Date:   startOfWeek.Add(time.Duration(i*24) * time.Hour),
			}

			_, err := uc.CreateEvent(ctx, createCmd)
			assert.NoError(t, err)
		}

		query := queries.GetEventsForWeek{
			UserID:    userID.String(),
			DateStart: startOfWeek,
		}

		events, err := uc.GetEventsForWeek(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, events, 3)
	})

	t.Run("events outside week range", func(t *testing.T) {
		userID := uuid.New()
		startOfWeek := time.Now().Truncate(24 * time.Hour)

		createCmd := commands.CreateEvent{
			UserID: userID.String(),
			Event:  "Next Week Event",
			Date:   startOfWeek.Add(8 * 24 * time.Hour),
		}

		_, err := uc.CreateEvent(ctx, createCmd)
		assert.NoError(t, err)

		query := queries.GetEventsForWeek{
			UserID:    userID.String(),
			DateStart: startOfWeek,
		}

		events, err := uc.GetEventsForWeek(ctx, query)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errorx.ErrNoEventsFound))
		assert.Empty(t, events)
	})
}

func TestCalendar_GetEventsForMonth(t *testing.T) {
	uc, ctx := setupTest()

	t.Run("successful get events for month", func(t *testing.T) {
		userID := uuid.New()
		startOfMonth := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		for i := 0; i < 5; i++ {
			createCmd := commands.CreateEvent{
				UserID: userID.String(),
				Event:  "Event " + string(rune('A'+i)),
				Date:   startOfMonth.Add(time.Duration(i*7*24) * time.Hour),
			}

			_, err := uc.CreateEvent(ctx, createCmd)
			assert.NoError(t, err)
		}

		query := queries.GetEventsForMonth{
			UserID:    userID.String(),
			DateStart: startOfMonth,
		}

		events, err := uc.GetEventsForMonth(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, events, 5)
	})

	t.Run("events outside month range", func(t *testing.T) {
		userID := uuid.New()
		startOfMonth := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		createCmd := commands.CreateEvent{
			UserID: userID.String(),
			Event:  "Next Month Event",
			Date:   startOfMonth.AddDate(0, 1, 1),
		}

		_, err := uc.CreateEvent(ctx, createCmd)
		assert.NoError(t, err)

		query := queries.GetEventsForMonth{
			UserID:    userID.String(),
			DateStart: startOfMonth,
		}

		events, err := uc.GetEventsForMonth(ctx, query)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errorx.ErrNoEventsFound))
		assert.Empty(t, events)
	})
}

func TestCalendar_parseUserID(t *testing.T) {
	uc, _ := setupTest()

	t.Run("valid UUID", func(t *testing.T) {
		validUUID := uuid.New().String()
		result, err := uc.parseUserID(validUUID, "test")

		assert.NoError(t, err)
		assert.Equal(t, validUUID, result.String())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		result, err := uc.parseUserID("invalid-uuid", "test")

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, result)
		assert.Contains(t, err.Error(), "invalid UUID")
	})
}

func TestCalendar_parseEventID(t *testing.T) {
	uc, _ := setupTest()

	t.Run("valid UUID", func(t *testing.T) {
		validUUID := uuid.New().String()
		result, err := uc.parseEventID(validUUID, "test")

		assert.NoError(t, err)
		assert.Equal(t, validUUID, result.String())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		result, err := uc.parseEventID("invalid-uuid", "test")

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, result)
		assert.Contains(t, err.Error(), "invalid UUID")
	})
}

func TestNewCalendar(t *testing.T) {
	storage := memory.NewCalendarStorage()
	repo := repository.NewCalendarRepo(storage)
	logger := slog.Default()

	uc := NewCalendar(logger, repo)

	assert.NotNil(t, uc)
	assert.Equal(t, logger, uc.log)
	assert.Equal(t, repo, uc.repo)
}
