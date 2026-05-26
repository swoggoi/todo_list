package users_transport_http

import (
	"net/http"

	"github.com/swoggoi/todo_list/internal/core/domain"
	core_logger "github.com/swoggoi/todo_list/internal/core/logger"
	core_http_request "github.com/swoggoi/todo_list/internal/core/transport/http/request"
	core_http_responce "github.com/swoggoi/todo_list/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"                    example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"   example:"+79998887766"`
}
type CreateUserResponse UserDTOResponse

// CreateUser 			godoc
// @Summary 			Создать пользователя
// @Description 		Создать нового пользователя в системе
// @Tags 				users
// @Accept 				json
// @Produce 			json
// @Param 				request body CreateUserRequest true "CreateUser тело запроса"
// @Success 			201 {object} CreateUserResponse "Успешно созданный пользователь"
// @Failure 			400 {object} core_http_responce.ErrorResponse "Bad Request"
// @Failure 			500 {object} core_http_responce.ErrorResponse "Internal server error"
// @Router 				/users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	log.Debug("invoke CreateUser handler")
	responseHandler := core_http_responce.NewHTTPResponseHandler(log, rw)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}
	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failded to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUnitialized(dto.FullName, dto.PhoneNumber)
}
