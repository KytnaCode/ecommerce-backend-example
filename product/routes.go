package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kytnacode/ecommerce-backend-example/validation"
)

// Routes implements routes for products using dependency injection
type Routes struct {
	productService Service
}

func NewRoutes(ps Service) *Routes {
	return &Routes{
		productService: ps,
	}
}

/*
GetProduct is an higher order function that takes an id validator for strings an return
the handler for a single product by the `id` parameter in path.
*/
func (rs *Routes) GetProduct(
	idValidator validation.Validator[string],
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate
		id := r.PathValue("id")
		if err := idValidator.Validate(&id); err != nil {
			http.Error(
				w,
				fmt.Sprintf("400: Bad Request: Invalid ID: %v", err),
				http.StatusBadRequest,
			)
			return
		}

		// Fetch
		err, p := rs.productService.Get(id)
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			http.Error(
				w,
				fmt.Sprintf("404: Not found product with id: %v", notFoundError.id),
				http.StatusNotFound,
			)
			return
		}

		// Send product
		enc := json.NewEncoder(w)

		err = enc.Encode(p)
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// GetAll is the handler for get all the products.
func (rs *Routes) GetAll(w http.ResponseWriter, r *http.Request) {
	// Fetch
	err, ps := rs.productService.GetAll()
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send
	enc := json.NewEncoder(w)
	if err = enc.Encode(ps); err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
}
