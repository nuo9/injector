package injector

import (
	"fmt"
	"reflect"
	"strings"
)

func Inject(obj interface{}, paths []string, value interface{}) error {
	field, e := travelPath(obj, paths)
	if e != nil {
		return e
	}

	setValue, e := createValue(field.Type(), value)
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

func createValue(required reflect.Type, value interface{}) (reflect.Value, error) {
	suppliedType := reflect.TypeOf(value)
	convertible := suppliedType.ConvertibleTo(required)
	if !convertible {
		return reflect.Value{}, fmt.Errorf("type %s is not assignable to %s", suppliedType, required)
	}

	return reflect.ValueOf(value).Convert(required), nil
}

func setFieldValue(field reflect.Value, set reflect.Value) error {
	field.Set(set)
	return nil
}
