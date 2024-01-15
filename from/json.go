package from

import (
	"encoding/json"
	"net/http"

	"github.com/bridgex-eu/wirex"
)

type JsonData[T any] struct {
	Data *T
}

func Json[T any](data *T) *JsonData[T] {
	return &JsonData[T]{Data: data}
}

func (j *JsonData[T]) FromRequest(r *http.Request) wirex.HTTPError {
	// TODO: add validation
	err := json.NewDecoder(r.Body).Decode(j.Data)
	if err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			return wirex.Error(http.StatusBadRequest, err)
		}
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return wirex.Error(http.StatusBadRequest, err)
		}

		return wirex.Error(http.StatusInternalServerError, err)
	}

	return nil
}
