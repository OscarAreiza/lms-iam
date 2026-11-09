// Package response provides the standard JSON envelope helpers used by every
// HTTP handler — see rules/2-anexos/C-api-hexagonal.md ("El contrato
// público" § Errores). Matches every sibling service's envelope shape
// exactly (error, message, traceId).
package response

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	TraceID string `json:"traceId"`
}

func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func Error(w http.ResponseWriter, status int, code, message, traceID string) {
	JSON(w, status, ErrorBody{Error: code, Message: message, TraceID: traceID})
}
