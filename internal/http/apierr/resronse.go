package apierr

import (
    "encoding/json"
    "net/http"
)

type APIError struct {
    Code       string            `json:"code"`
    Message    string            `json:"message"`
    StatusCode int               `json:"-"`
    Fields     map[string]string `json:"fields,omitempty"`
}

func (e *APIError) Error() string {
    return e.Message
}

func NewValidationError(fields map[string]string) *APIError {
    return &APIError{
        Code:       "validation_error",
        Message:    "Некорректные данные",
        StatusCode: http.StatusBadRequest,
        Fields:     fields,
    }
}

func NewInternalError(msg string) *APIError {
    return &APIError{
        Code:       "internal_error",
        Message:    msg,
        StatusCode: http.StatusInternalServerError,
    }
}

func RespondWithError(w http.ResponseWriter, err *APIError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(err.StatusCode)
    json.NewEncoder(w).Encode(err)
}
