package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"task/domain/task/service"
	"task/internal/handler"
	"task/internal/logging"
)

const (
	url = "/api/task/:id"
)

type Handler struct {
	logger  *logging.Logger
	service service.Service
}

func NewHandler(logger *logging.Logger, service service.Service) *Handler {

	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) Register(router handler.HandlerFunc) {
	router.HandlerFunc(http.MethodGet, url, h.GetTaskInfo)
}

func (h *Handler) GetTaskInfo(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("id")

	taskId, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("invalid id", "id", id)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("invalid id"))
		if err != nil {
			h.logger.ErrorErr(err)
			return
		}
		return
	}

	h.logger.Info("get task", "id", taskId)

	task := h.service.GetOrCreate(ctx, taskId)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(task.String()))
	if err != nil {
		h.logger.ErrorErr(err)
		return
	}

}
