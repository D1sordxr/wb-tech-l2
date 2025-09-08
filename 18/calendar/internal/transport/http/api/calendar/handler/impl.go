package handler

import (
	"context"
	"errors"
	"fmt"
	"wb-tech-l2/18/calendar/internal/application/calendar/command"
	"wb-tech-l2/18/calendar/internal/application/calendar/ports"
	"wb-tech-l2/18/calendar/internal/application/calendar/query"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/errorx"
	"wb-tech-l2/18/calendar/internal/domain/core/calendar/model"
	"wb-tech-l2/18/calendar/pkg/httputil"

	"github.com/oapi-codegen/runtime/types"
)

type Handlers struct {
	uc ports.CalendarUseCase
}

func NewHandlers(uc ports.CalendarUseCase) *Handlers {
	return &Handlers{uc: uc}
}

func (h Handlers) PostCreateEvent(ctx context.Context, request PostCreateEventRequestObject) (PostCreateEventResponseObject, error) {
	t, err := httputil.ParseGenTime(request.JSONBody.Date)
	if err != nil {
		return PostCreateEvent400JSONResponse{
			Error: "invalid date format",
		}, nil
	}

	eventID, err := h.uc.CreateEvent(ctx, command.CreateEvent{
		UserID: request.JSONBody.UserId,
		Date:   t,
		Event:  request.JSONBody.Event,
	})
	if err != nil {
		return PostCreateEvent503JSONResponse(ErrorResponse{
			Error: err.Error(),
		}), nil
	}

	return PostCreateEvent200JSONResponse(SuccessResponse{
		Result: fmt.Sprintf("Event created with ID: %s", eventID),
	}), nil
}

func (h Handlers) PostUpdateEvent(ctx context.Context, request PostUpdateEventRequestObject) (PostUpdateEventResponseObject, error) {
	t, err := httputil.ParseGenTime(request.JSONBody.Date)
	if err != nil {
		return PostUpdateEvent400JSONResponse{
			Error: "invalid date format",
		}, nil
	}

	if err = h.uc.UpdateEvent(ctx, command.UpdateEvent{
		ID:     request.JSONBody.Id,
		UserID: request.JSONBody.UserId,
		Date:   t,
		Event:  request.JSONBody.Event,
	}); err != nil {
		return PostUpdateEvent503JSONResponse(ErrorResponse{
			Error: err.Error(),
		}), nil
	}

	return PostUpdateEvent200JSONResponse(SuccessResponse{
		Result: fmt.Sprintf("Event %s updated successfully", request.JSONBody.Id),
	}), nil
}

func (h Handlers) PostDeleteEvent(ctx context.Context, request PostDeleteEventRequestObject) (PostDeleteEventResponseObject, error) {
	if err := h.uc.DeleteEvent(ctx, command.DeleteEvent{
		ID:     request.JSONBody.Id,
		UserID: request.JSONBody.UserId,
	}); err != nil {
		return PostDeleteEvent503JSONResponse(ErrorResponse{
			Error: err.Error(),
		}), nil
	}

	return PostDeleteEvent200JSONResponse(SuccessResponse{
		Result: fmt.Sprintf("Event %s deleted successfully", request.JSONBody.Id),
	}), nil
}

func (h Handlers) GetEventsForDay(ctx context.Context, request GetEventsForDayRequestObject) (GetEventsForDayResponseObject, error) {
	t, err := httputil.ParseGenDateOnly(request.Params.Date)
	if err != nil {
		return GetEventsForDay400JSONResponse{
			Error: "invalid date format",
		}, nil
	}

	events, err := h.uc.GetEventsForDay(ctx, query.GetEventsForDay{
		UserID: request.Params.UserId,
		Date:   t,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrNoEventsFound):
			return GetEventsForDay400JSONResponse(ErrorResponse{Error: err.Error()}), nil
		default:
			return GetEventsForDay500Response{}, nil
		}
	}

	resEvents := make([]CalendarEvent, len(events))
	for i, event := range events {
		resEvents[i] = h.parseCalendarEventFromModel(event)
	}

	return GetEventsForDay200JSONResponse(EventsListResponse{
		Result: resEvents,
	}), nil
}

func (h Handlers) GetEventsForWeek(ctx context.Context, request GetEventsForWeekRequestObject) (GetEventsForWeekResponseObject, error) {
	tStart, err := httputil.ParseGenDateOnly(request.Params.DateStart)
	if err != nil {
		return GetEventsForWeek400JSONResponse(ErrorResponse{
			Error: "invalid date format",
		}), nil
	}

	events, err := h.uc.GetEventsForWeek(ctx, query.GetEventsForWeek{
		UserID:    request.Params.UserId,
		DateStart: tStart,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrNoEventsFound):
			return GetEventsForWeek400JSONResponse(ErrorResponse{Error: err.Error()}), nil
		default:
			return GetEventsForWeek500Response{}, nil
		}
	}

	resEvents := make([]CalendarEvent, len(events))
	for i, event := range events {
		resEvents[i] = h.parseCalendarEventFromModel(event)
	}

	return GetEventsForWeek200JSONResponse(EventsListResponse{
		Result: resEvents,
	}), nil
}

func (h Handlers) GetEventsForMonth(ctx context.Context, request GetEventsForMonthRequestObject) (GetEventsForMonthResponseObject, error) {
	tStart, err := httputil.ParseGenDateOnly(request.Params.DateStart)
	if err != nil {
		return GetEventsForMonth400JSONResponse(ErrorResponse{
			Error: "invalid date format",
		}), nil
	}

	events, err := h.uc.GetEventsForMonth(ctx, query.GetEventsForMonth{
		UserID:    request.Params.UserId,
		DateStart: tStart,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrNoEventsFound):
			return GetEventsForMonth400JSONResponse(ErrorResponse{Error: err.Error()}), nil
		default:
			return GetEventsForMonth500Response{}, nil
		}
	}

	resEvents := make([]CalendarEvent, len(events))
	for i, event := range events {
		resEvents[i] = h.parseCalendarEventFromModel(event)
	}

	return GetEventsForMonth200JSONResponse(EventsListResponse{
		Result: resEvents,
	}), nil
}

func (Handlers) parseCalendarEventFromModel(event *model.CalendarEvent) CalendarEvent {
	eId := event.ID.String()
	uId := event.UserID.String()
	return CalendarEvent{
		CreatedAt: &event.CreatedAt,
		Date:      &types.Date{Time: event.Date},
		Event:     &event.Event,
		Id:        &eId,
		UpdatedAt: &event.UpdatedAt,
		UserId:    &uId,
	}
}
