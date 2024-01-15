package from

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bridgex-eu/wirex"
)

type QueryData[T decodable] struct {
	Data    *T
	Name    string
	Default *T
}

func (q *QueryData[T]) FromRequest(r *http.Request) wirex.HTTPError {
	val := r.URL.Query().Get(q.Name)
	if val == "" {
		if q.Default == nil {
			return wirex.Error(http.StatusBadRequest, errors.New("query parameter: "+q.Name+" not found"))
		}

		*q.Data = *q.Default
		return nil
	}

	decoded, err := decode[T](val)
	if err != nil {
		return wirex.Error(http.StatusBadRequest, fmt.Errorf("query parameter: %s has wrong type, value: %s", q.Name, val))
	}

	*q.Data = decoded
	return nil
}

func Query[T decodable](name string, value *T) *QueryData[T] {
	return &QueryData[T]{Data: value, Name: name, Default: nil}
}

func QueryOr[T decodable](name string, value *T, defaultValue T) *QueryData[T] {
	return &QueryData[T]{Data: value, Name: name, Default: &defaultValue}
}
