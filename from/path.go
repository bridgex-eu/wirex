package from

import (
	"fmt"
	"net/http"

	"github.com/bridgex-eu/wirex"
)

type PathData[T decodable] struct {
	Data *T
	Name string
}

func (p *PathData[T]) FromRequest(r *http.Request) wirex.HTTPError {
	val := r.PathValue(p.Name)
	if val == "" {
		return wirex.Error(http.StatusBadRequest, fmt.Errorf("path parameter: %s not found", p.Name))
	}

	decoded, err := decode[T](val)
	if err != nil {
		return wirex.Error(http.StatusBadRequest, fmt.Errorf("path parameter: %s has wrong type, value: %s", p.Name, val))
	}

	*p.Data = decoded
	return nil
}

func Path[T decodable](name string, value *T) *PathData[T] {
	return &PathData[T]{Data: value, Name: name}
}
