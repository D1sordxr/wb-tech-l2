package memory

import (
	"context"
	"sync"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/model"
)

type CalendarStorage struct {
	storage map[string]map[string]*model.CalendarEvent // userID -> eventID -> *CalendarEvent
	mutex   sync.RWMutex
}

func NewCalendarStorage() *CalendarStorage {
	return &CalendarStorage{
		storage: make(map[string]map[string]*model.CalendarEvent),
		mutex:   sync.RWMutex{},
	}
}

func (c *CalendarStorage) Set(_ context.Context, userID string, event *model.CalendarEvent) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.storage[userID]; !ok {
		c.storage[userID] = make(map[string]*model.CalendarEvent)
	}

	c.storage[userID][event.ID.String()] = event
	return nil
}

func (c *CalendarStorage) Get(_ context.Context, userID, eventID string) (*model.CalendarEvent, error) {
	c.mutex.RLock()
	userEvents, ok := c.storage[userID]
	c.mutex.RUnlock()
	if !ok {
		return nil, ErrEventDoesNotExist
	}

	event, ok := userEvents[eventID]
	if !ok {
		return nil, ErrEventDoesNotExist
	}

	return event, nil
}

func (c *CalendarStorage) Exists(_ context.Context, userID, eventID string) bool {
	c.mutex.RLock()
	userEvents, ok := c.storage[userID]
	c.mutex.RUnlock()
	if !ok {
		return false
	}

	if _, ok = userEvents[eventID]; !ok {
		return false
	}

	return true
}

func (c *CalendarStorage) List(_ context.Context, userID string) ([]*model.CalendarEvent, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	userEvents, ok := c.storage[userID]
	if !ok {
		return []*model.CalendarEvent{}, nil
	}

	events := make([]*model.CalendarEvent, 0, len(userEvents))
	for _, event := range userEvents {
		events = append(events, event)
	}

	return events, nil
}

func (c *CalendarStorage) Delete(_ context.Context, userID, eventID string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	userEvents, ok := c.storage[userID]
	if !ok {
		return ErrEventDoesNotExist
	}
	delete(userEvents, eventID)
	return nil
}
