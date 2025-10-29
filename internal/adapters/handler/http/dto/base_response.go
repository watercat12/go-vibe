package dto

import "net/http"

type Response struct {
	Status            int    `json:"status" example:"200"`
	Message           string `json:"message" example:"Success"`
	Data              any    `json:"data"`
}

var (
	SuccessResponse = Response{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    nil,
	}
	InternalErrorResponse = Response{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Data:    nil,
	}
	BadRequestResponse = Response{
		Status:  http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Data:    nil,
	}
	UnauthorizedResponse = Response{
		Status:  http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
		Data:    nil,
	}
)
