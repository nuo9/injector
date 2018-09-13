package injector

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Inject(obj interface{}, paths []string, value interface{}) error {
	field, e := travelPath(obj, paths)
	if e != nil {
		return e
	}

	setValue, e := convertValue(field.Type(), value)
	if e != nil {
		return e
	}

	return setFieldValue(field, setValue)
}

func travelPath(obj interface{}, paths []string) (reflect.Value, error) {
	v := reflect.ValueOf(obj).Elem()

	for _, f := range paths {
		v = v.FieldByName(f)
		if !v.IsValid() {
			return reflect.Value{}, fmt.Errorf("field %s is not valid", f)
		}
	}

	if !v.CanSet() {
		return reflect.Value{}, fmt.Errorf("fields %s cannot be set", strings.Join(paths, "."))
	}

	return v, nil
}

func convertValue(required reflect.Type, value interface{}) (reflect.Value, error) {
	v, e := convertFromString(required, value)
	if e == nil {
		return v, nil
	}

	v, e = convertToString(required, value)
	if e == nil {
		return v, nil
	}

	return convertByReflect(required, value)
}

func convertByReflect(required reflect.Type, value interface{}) (reflect.Value, error) {
	suppliedType := reflect.TypeOf(value)
	convertible := suppliedType.ConvertibleTo(required)
	if !convertible {
		return reflect.Value{}, fmt.Errorf("type %s is not convertible to %s", suppliedType, required)
	}

	return reflect.ValueOf(value).Convert(required), nil
}

func convertFromString(required reflect.Type, value interface{}) (reflect.Value, error) {
	s, ok := value.(string)
	if !ok {
		return reflect.Value{}, fmt.Errorf("type %s is not string", reflect.TypeOf(value))
	}

	switch required.Kind() {
	case reflect.Bool:
		b, e := strconv.ParseBool(s)
		return reflect.ValueOf(b), e
	case reflect.Int:
		i, e := strconv.ParseInt(s, 10, 0)
		return reflect.ValueOf(int(i)), e
	case reflect.Int8:
		i, e := strconv.ParseInt(s, 10, 8)
		return reflect.ValueOf(int8(i)), e
	case reflect.Int16:
		i, e := strconv.ParseInt(s, 10, 16)
		return reflect.ValueOf(int16(i)), e
	case reflect.Int32:
		i, e := strconv.ParseInt(s, 10, 32)
		return reflect.ValueOf(int32(i)), e
	case reflect.Int64:
		i, e := strconv.ParseInt(s, 10, 64)
		return reflect.ValueOf(int64(i)), e
	case reflect.Uint:
		u, e := strconv.ParseUint(s, 10, 0)
		return reflect.ValueOf(uint(u)), e
	case reflect.Uint8:
		u, e := strconv.ParseUint(s, 10, 8)
		return reflect.ValueOf(uint8(u)), e
	case reflect.Uint16:
		u, e := strconv.ParseUint(s, 10, 16)
		return reflect.ValueOf(uint16(u)), e
	case reflect.Uint32:
		u, e := strconv.ParseUint(s, 10, 32)
		return reflect.ValueOf(uint32(u)), e
	case reflect.Uint64:
		u, e := strconv.ParseUint(s, 10, 64)
		return reflect.ValueOf(uint64(u)), e
	case reflect.Float32:
		f, e := strconv.ParseFloat(s, 32)
		return reflect.ValueOf(float32(f)), e
	case reflect.Float64:
		f, e := strconv.ParseFloat(s, 64)
		return reflect.ValueOf(float64(f)), e
	default:
		return reflect.Value{}, fmt.Errorf("type string cannot convert to %s", required)
	}
}

func convertToString(required reflect.Type, value interface{}) (reflect.Value, error) {
	if !(required.Kind() == reflect.String) {
		return reflect.Value{}, fmt.Errorf("required type %s is not string", required)
	}

	return reflect.ValueOf(fmt.Sprint(value)), nil
}

func setFieldValue(field reflect.Value, set reflect.Value) error {
	field.Set(set)
	return nil
}
