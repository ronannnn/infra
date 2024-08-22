package errors_test

import (
	"fmt"
	"testing"

	"github.com/ronannnn/infra/errors"
)

func TestNew(t *testing.T) {
	e := A()
	fmt.Printf("%s\n\n", e)
	fmt.Printf("%v\n\n", e)
}

func A() error {
	return B()
}

func B() error {
	return C()
}

func C() error {
	return D()
}

func D() error {
	return errors.New("internal server error").
		WithMsg("db connection error!!").
		WithError(fmt.Errorf("db connection error")).
		WithStack()
}
