package wirex

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type HTTPError interface {
	Writer
	error
}

type DefaultHTTPError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

var _ HTTPError = &DefaultHTTPError{}

func (e *DefaultHTTPError) WriteResponse(w http.ResponseWriter, r *http.Request) {
	writeHeader := w.Header()
	if writeHeader.Get(HeaderContentType) == "" {
		writeHeader.Set(HeaderContentType, MIMEApplicationJSON)
	}

	w.WriteHeader(e.Status)

	if err := json.NewEncoder(w).Encode(*e); err != nil {
		slog.Error("cannot write json to response", "error", err)
	}
}

func (s *DefaultHTTPError) Error() string {
	return s.Message
}

func Error(status int, err error) HTTPError {
	return &DefaultHTTPError{status, err.Error()}
}
