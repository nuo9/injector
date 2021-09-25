package injector

import (
	"testing"
)

type (
	S1 struct {
		Id string
	}

	S2 struct {
		S1 S1
	}

	S3 struct {
		S2 S2
	}

	S4 struct {
		S3 S3
	}
)

var (
	testString = "this is a test string"

	s1 = S1{}
	s3 = S4{S3: S3{
		S2: S2{
			S1: S1{},
		},
	}}
)

func BenchmarkInjectorNesting1(b *testing.B) {
	e := Inject(&s1, []string{"Id"}, testString)
	if e != nil {
		panic(e)
	}
}

func BenchmarkInjectorNesting4(b *testing.B) {
	e := Inject(&s3, []string{"S3", "S2", "S1", "Id"}, testString)
	if e != nil {
		panic(e)
	}
}

func BenchmarkInjectConvertTypes(b *testing.B) {
	e := Inject(&s3, []string{"S3", "S2", "S1", "Id"}, 1000.0)
	if e != nil {
		panic(e)
	}
}
