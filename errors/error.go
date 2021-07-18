package errors

import (
	"net/http"
)

type AppointmentErr interface {
	GetMessage() string
	GetStatus() int
	GetError() string
}

type appointmentErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func (e *appointmentErr) GetError() string {
	return e.Error
}

func (e *appointmentErr) GetMessage() string {
	return e.Message
}

func (e *appointmentErr) GetStatus() int {
	return e.Status
}

func NewNotFoundError(message string, err error) AppointmentErr {
	return &appointmentErr{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   err.Error(),
	}
}

func NewBadRequestError(message string, err error) AppointmentErr {
	return &appointmentErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   err.Error(),
	}
}

func NewUnprocessibleEntityError(message string, err error) AppointmentErr {
	return &appointmentErr{
		Message: message,
		Status:  http.StatusUnprocessableEntity,
		Error:   err.Error(),
	}
}

func NewInternalServerError(message string, err error) AppointmentErr {
	return &appointmentErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   err.Error(),
	}
}

func NewGeneralError(message string, err error) AppointmentErr {
	errMsg := ""

	if err != nil {
		errMsg = err.Error()
	}

	return &appointmentErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   errMsg,
	}
}

func NewGeneralForbiddenError(message string, err error) AppointmentErr {
	errMsg := ""

	if err != nil {
		errMsg = err.Error()
	}

	return &appointmentErr{
		Message: message,
		Status:  http.StatusForbidden,
		Error:   errMsg,
	}
}
