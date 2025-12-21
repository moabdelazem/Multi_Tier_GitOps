package pkg

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func JSONSuccess(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: message})
}

func NotFound(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: message})
}

func InternalError(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: message})
}

func ServiceUnavailable(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusServiceUnavailable, data)
}
