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

type CreateRequest struct{}

type CreateResponse struct {
	ID string `json:"id"`
}

func FromCreateRequest(req CreateRequest) Product {
	return Product{}
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

/*
Create handles product create requests

It uses a validator of `CreateRequest` to enforce request body validation rules.
*/
func (rs *Routes) Create(reqValidator validation.Validator[CreateRequest]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Read body
		var req *CreateRequest
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&req)
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Validate
		if err := reqValidator.Validate(req); err != nil {
			http.Error(w, fmt.Sprintf("400: Bad Request: %v", err), http.StatusBadRequest)
			return
		}

		// Create
		p := FromCreateRequest(*req)

		if err = rs.productService.Create(&p); err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Respond
		res := CreateResponse{ID: p.ID}

		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err = enc.Encode(res); err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
