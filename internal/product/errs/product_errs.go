package errs

import (
	"errors"
	"fmt"
)

var (
	ProductErrorInvalidUUID      = errors.New("Invalid uuid")
	ProductErrorCategoryNotFound = errors.New("Product category not found")
)

type ProductErrs struct {
	Err error
}

func (e ProductErrs) Error() error {
	return fmt.Errorf(e.Err.Error())
}
