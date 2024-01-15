package from

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/bridgex-eu/wirex"
)

type FormData[T any] struct {
	Data *T
}

func Form[T any](data *T) *FormData[T] {
	return &FormData[T]{Data: data}
}

func (f *FormData[T]) FromRequest(r *http.Request) wirex.HTTPError {
	// Parse the form data from the request
	err := r.ParseForm()
	if err != nil {
		return wirex.Error(http.StatusBadRequest, err)
	}

	err = decodeForm(r.Form, f.Data)
	if err != nil {
		if _, ok := err.(*url.EscapeError); ok {
			return wirex.Error(http.StatusBadRequest, err)
		}

		return wirex.Error(http.StatusInternalServerError, err)
	}

	return nil
}

func decodeForm(form url.Values, to any) error {
	toValue := reflect.ValueOf(to)

	// Check if the 'to' parameter is a pointer to a struct
	if toValue.Kind() != reflect.Ptr || toValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("the 'to' argument must be a pointer to a struct")
	}

	structValue := toValue.Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Get the form field name from the struct tag (or use the struct field name)
		formFieldName := fieldType.Tag.Get("form")
		if formFieldName == "" {
			formFieldName = fieldType.Name
		}

		// Check if the form has this field
		if values, ok := form[formFieldName]; ok && len(values) > 0 {
			err := decodeValue(values[0], field.Addr().Interface())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
