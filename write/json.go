package write

import (
	"encoding/json"
	"net/http"

	"github.com/bridgex-eu/wirex"
)

type JsonData[T any] struct {
	Status int
	Data   *T
}

var _ wirex.Writer = &JsonData[int]{}

func (j *JsonData[T]) WriteResponse(w http.ResponseWriter, r *http.Request) {
	SetContentType(w, wirex.MIMEApplicationJSON)

	w.WriteHeader(j.Status)

	if err := json.NewEncoder(w).Encode(j.Data); err != nil {
		String(http.StatusInternalServerError, err.Error()).WriteResponse(w, r)
	}
}

func Json[T any](status int, data *T, other ...wirex.HeaderWriter) wirex.Writer {
	return WithHeader(&JsonData[T]{status, data}, other...)
}
