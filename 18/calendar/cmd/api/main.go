package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"wb-tech-l2/18/calendar/internal/application/calendar/usecase"
	"wb-tech-l2/18/calendar/internal/infrastucture/config"
	"wb-tech-l2/18/calendar/internal/infrastucture/storage/calendar/memory"
	"wb-tech-l2/18/calendar/internal/infrastucture/storage/calendar/memory/repository"
	"wb-tech-l2/18/calendar/internal/transport/http"
	"wb-tech-l2/18/calendar/internal/transport/http/api/calendar"
	"wb-tech-l2/18/calendar/internal/transport/http/api/calendar/handler"

	loadApp "wb-tech-l2/18/calendar/internal/infrastucture/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := config.NewConfig()
	log := slog.Default()

	calendarStorage := memory.NewCalendarStorage()
	calendarRepo := repository.NewCalendarRepo(calendarStorage)
	calendarUseCase := usecase.NewCalendar(log, calendarRepo)
	calendarHandlers := handler.NewHandlers(calendarUseCase)
	calendarRouteRegisterer := calendar.NewRouteRegisterer(calendarHandlers)

	server := http.NewServer(
		log,
		&cfg.Server,
		calendarRouteRegisterer,
	)

	app := loadApp.NewApp(log, server)
	app.Run(ctx)
}
