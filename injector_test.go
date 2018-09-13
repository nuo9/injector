package injector

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)

type simpleAnimal struct {
	Id   uint32
	Name string
	Bird struct {
		Flying bool
	}
	Dog struct {
		Husky struct {
			IQ uint8
		}
	}
}

const content string = `
Id: 10,
Name: "2333",
Bird: {
 Flying: true
},
Dog: {
 Husky: {
  IQ: 20
 }
}
`

func TestInject(t *testing.T) {
	d := simpleAnimal{}
	json.Unmarshal([]byte(content), &d)

	errors := make([]error, 0)
	var e error

	e = validate(&d, []string{"Id"}, 666, func(d *simpleAnimal) interface{} {
		return d.Id
	})
	errors = append(errors, e)

	e = validate(&d, []string{"Name"}, "haha", func(d *simpleAnimal) interface{} {
		return d.Name
	})
	errors = append(errors, e)

	e = validate(&d, []string{"Name"}, 600, func(d *simpleAnimal) interface{} {
		return d.Name
	})
	errors = append(errors, e)

	e = validate(&d, []string{"Bird", "Flying"}, false, func(d *simpleAnimal) interface{} {
		return d.Bird.Flying
	})
	errors = append(errors, e)

	e = validate(&d, []string{"Dog", "Husky", "IQ"}, 11, func(d *simpleAnimal) interface{} {
		return d.Dog.Husky.IQ
	})
	errors = append(errors, e)

	e = validate(&d, []string{"Dog", "Husky", "IQ"}, "12", func(d *simpleAnimal) interface{} {
		return d.Dog.Husky.IQ
	})
	errors = append(errors, e)

	e = validate(&d, []string{"Dog"},
		struct {
			Husky struct {
				IQ uint8
			}
		}{Husky: struct{ IQ uint8 }{IQ: 7}},
		func(d *simpleAnimal) interface{} {
			return d.Dog
		})
	errors = append(errors, e)

	hasError := false
	for _, e := range errors {
		if e != nil {
			log.Println(e)
			hasError = true
		}
	}
	if hasError {
		t.Error("has errors")
	}
}

func validate(d *simpleAnimal, paths []string, value interface{}, getValue func(d *simpleAnimal) interface{}) error {
	e := Inject(d, paths, value)
	if e != nil {
		return e
	}
	if !equals(value, getValue(d)) {
		return fmt.Errorf("field %s not setted", strings.Join(paths, "."))
	}
	return nil
}

func equals(a, b interface{}) bool {
	as := reflect.TypeOf(a).Kind() == reflect.String
	bs := reflect.TypeOf(b).Kind() == reflect.String
	if as || bs {
		return fmt.Sprint(a) == fmt.Sprint(b)
	}

	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	a2b := av.Type().ConvertibleTo(bv.Type())
	if a2b {
		av = av.Convert(bv.Type())
		return reflect.DeepEqual(av.Interface(), bv.Interface())
	}

	b2a := bv.Type().ConvertibleTo(av.Type())
	if b2a {
		bv = bv.Convert(av.Type())
		return reflect.DeepEqual(av.Interface(), bv.Interface())
	}

	return a == b
}
