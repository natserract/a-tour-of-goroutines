package errs

import (
	"errors"
	"fmt"
)

var (
	ProductErrsCategoryNotFound = errors.New("Product category not found")
	ProductErrsSkuOverflow      = errors.New("Product sku is overflow")
	ProductErrsImageUrlInvalid  = errors.New("Product image url invalid")
)

type ProductErrs struct {
	Err error
}

func (e ProductErrs) Error() error {
	return fmt.Errorf(e.Err.Error())
}
