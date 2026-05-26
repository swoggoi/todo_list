package tasks_transport

import (
	"net/http"

	core_logger "github.com/swoggoi/todo_list/internal/core/logger"
	core_http_request "github.com/swoggoi/todo_list/internal/core/transport/http/request"
	core_http_responce "github.com/swoggoi/todo_list/internal/core/transport/http/response"
)

// DeleteTask godoc
// @Summary Удаление задачи
// @Description Удалить существующую в системе задачу по её ID
// @Tags tasks
// @Param id path int true "ID удаляемой задачи"
// @Success 204 "Успешное удаление задачи"
// @Failure 400 {object} core_http_responce.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_responce.ErrorResponse "task not found"
// @Failure 500 {object} core_http_responce.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_responce.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}
	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.NoContentResponse()
}
