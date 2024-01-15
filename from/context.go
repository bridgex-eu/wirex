package from

import (
	"fmt"
	"net/http"

	"github.com/bridgex-eu/wirex"
)

type ContextData[T any] struct {
	Data *T
	Key  string
}

func (c *ContextData[T]) FromRequest(r *http.Request) wirex.HTTPError {
	val := r.Context().Value(c.Key)
	if val == nil {
		return wirex.Error(http.StatusInternalServerError, fmt.Errorf("context value for key: %s not found", c.Key))
	}

	typed, ok := val.(T)
	if !ok {
		return wirex.Error(http.StatusInternalServerError, fmt.Errorf("context value for key: %s is not of the expected type", c.Key))
	}

	*c.Data = typed
	return nil
}

func FromContext[T any](key string, value *T) *ContextData[T] {
	return &ContextData[T]{Data: value, Key: key}
}
