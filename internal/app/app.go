package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"task/domain/task/handler"
	"task/domain/task/service"
	"task/domain/task/storage"
	"task/internal/config"
	"task/internal/logging"
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *httprouter.Router
	httpServer *http.Server

	taskService service.Service
}

func NewApp(cfg *config.Config, logger *logging.Logger) *App {

	logger.Info("router init")
	router := httprouter.New()

	taskStorage := storage.NewStorage(logger)
	taskService := service.NewService(logger, taskStorage)
	taskHandler := handler.NewHandler(logger, taskService)
	taskHandler.Register(router)

	return &App{
		cfg:         cfg,
		logger:      logger,
		router:      router,
		taskService: taskService,
	}
}

func (a *App) Run() {

	a.logger.Info("server start", "host", a.cfg.Host, "port", a.cfg.Port)
	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port),
		Handler: a.router,
	}

	if err := a.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		a.logger.ErrorErr(err)
		return
	}

}
func (a *App) Stop(ctx context.Context) {

	tasksInProgress := a.taskService.CountTasksInProgress(ctx)
	fmt.Printf("tasks in progress: %d\n", tasksInProgress)

	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		a.logger.ErrorErr(err)
	}
	a.logger.Info("server stop")

	tasksInfo := a.taskService.Shutdown(ctx)
	fmt.Printf("tasks info:\n%v", tasksInfo)
}
