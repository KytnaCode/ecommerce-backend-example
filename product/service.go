package product

import "fmt"

// NotFoundError is an error for not found products by a given id.
type NotFoundError struct {
	id string
}

// ID returns the not found id.
func (e NotFoundError) ID() string {
	return e.id
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("could not found product with id: %v", e.id)
}

// Product is a model for a product.
type Product struct{}

// Service abstracts product access methods.
type Service interface {
	/*
		Get return a product with the given id.

		If product is found succefully Get returns a nil error and a pointer to the product, if
		the product is not found, returns a `NotFoundError` and a nil pointer, if a error ocurr
		on Get, then a non-nil error and a nil pointer will be returned.
	*/
	Get(id string) (error, *Product)

	// GetAll returns all products in the service.
	GetAll() (error, []Product)
}
