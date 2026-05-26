package tasks_transport

import (
	"net/http"

	"github.com/swoggoi/todo_list/internal/core/domain"
	core_logger "github.com/swoggoi/todo_list/internal/core/logger"
	core_http_request "github.com/swoggoi/todo_list/internal/core/transport/http/request"
	core_http_responce "github.com/swoggoi/todo_list/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"Поиграть"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"С друзьями в пк клубе в 16:00"`
	AuthorUserID int     `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask			godoc
// @Summary 			Создать задачу
// @Description 		Создать новую задачу в системе
// @Tags 				tasks
// @Accept 				json
// @Produce 			json
// @Param 				request body CreateTaskRequest true "CreateTask тело запроса"
// @Success 			201 {object} CreateTaskResponse "Успешно созданная задача"
// @Failure 			400 {object} core_http_responce.ErrorResponse "Bad Request"
// @Failure 			404 {object} core_http_responce.ErrorResponse "Author not found"
// @Failure 			500 {object} core_http_responce.ErrorResponse "Internal server error"
// @Router 				/tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_responce.NewHTTPResponseHandler(log, rw)
	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	taskDomain := domain.NewTaskUnitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}
	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
