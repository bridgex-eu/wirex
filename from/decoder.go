package from

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

type decodable interface {
	bool | ~string | constraints.Integer | uuid.UUID
}

func decode[T decodable](val string) (T, error) {
	var res T

	switch any(res).(type) {
	case int:
		parsed, err := strconv.Atoi(val)
		return any(parsed).(T), err
	case int8:
		parsed, err := strconv.ParseInt(val, 10, 8)
		return any(int8(parsed)).(T), err
	case int16:
		parsed, err := strconv.ParseInt(val, 10, 16)
		return any(int16(parsed)).(T), err
	case int32:
		parsed, err := strconv.ParseInt(val, 10, 32)
		return any(int32(parsed)).(T), err
	case int64:
		parsed, err := strconv.ParseInt(val, 10, 64)
		return any(parsed).(T), err
	case uint:
		parsed, err := strconv.ParseUint(val, 10, 32)
		return any(uint(parsed)).(T), err
	case uint8:
		parsed, err := strconv.ParseUint(val, 10, 8)
		return any(uint8(parsed)).(T), err
	case uint16:
		parsed, err := strconv.ParseUint(val, 10, 16)
		return any(uint16(parsed)).(T), err
	case uint32:
		parsed, err := strconv.ParseUint(val, 10, 32)
		return any(uint32(parsed)).(T), err
	case uint64:
		parsed, err := strconv.ParseUint(val, 10, 64)
		return any(parsed).(T), err
	case bool:
		parsed, err := strconv.ParseBool(val)
		return any(parsed).(T), err
	case string:
		return any(val).(T), nil
	case uuid.UUID:
		parsed, err := uuid.Parse(val)
		return any(parsed).(T), err
	default:
		return res, fmt.Errorf("unsupported type: %s", reflect.TypeOf(res).String())
	}
}

func decodeValue(from string, to any) error {
	toValue := reflect.ValueOf(to)

	// Check if the 'to' parameter is a pointer
	if toValue.Kind() != reflect.Ptr || toValue.IsNil() {
		return fmt.Errorf("the 'to' argument must be a non-nil pointer")
	}

	// Dereference the pointer
	elem := toValue.Elem()

	switch elem.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(from, 10, 64)
		if err != nil {
			return err
		}
		elem.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(from, 10, 64)
		if err != nil {
			return err
		}
		elem.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(from, 64)
		if err != nil {
			return err
		}
		elem.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(from)
		if err != nil {
			return err
		}
		elem.SetBool(boolVal)
	case reflect.String:
		elem.SetString(from)
	case reflect.Struct:
		if _, ok := to.(*uuid.UUID); ok {
			parsedUUID, err := uuid.Parse(from)
			if err != nil {
				return err
			}
			elem.Set(reflect.ValueOf(parsedUUID))
		} else {
			return fmt.Errorf("struct type not supported: %s", elem.Type())
		}
	default:
		return fmt.Errorf("unsupported type: %s", elem.Type())
	}

	return nil
}