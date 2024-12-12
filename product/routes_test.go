package product_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kytnacode/ecommerce-backend-example/product"
)

type MockProductService struct {
	// getFn is used as body of MockProductService.Get to allow using the same mock for all tests.
	getFn func(id string) (error, *product.Product)
}

// Get implements `product.Service.Get`.
func (ps *MockProductService) Get(
	id string,
) (error, *product.Product) {
	return ps.getFn(id)
}

// AlwaysValid is a validator that always return a nil error.
type AlwaysValid[T any] struct{}

func (v AlwaysValid[T]) Validate(value *T) error {
	return nil
}

// AlwaysInvalid is a validator that always return a non-nil error.
type AlwaysInvalid[T any] struct{}

func (v AlwaysInvalid[T]) Validate(value *T) error {
	return errors.New("invalid value")
}

func TestGetProduct_shouldReturnOkResponse(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()

	rs := product.NewRoutes(&MockProductService{getFn: func(id string) (error, *product.Product) {
		return nil, &product.Product{} // Return an empty product.
	}})

	// Is required to register handler in a http.ServeMux instead of just call ServeHTTP on handler
	// to make path params work.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}", rs.GetProduct(AlwaysValid[string]{}))

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(
			"GetProduct should return an 200 status code: got: %v msg: %v",
			rr.Code,
			rr.Body.String(),
		)
	}
}

func TestGetProduct_shouldReturnNotFoundIfServiceReturnNotFoundError(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()

	rs := product.NewRoutes(&MockProductService{getFn: func(id string) (error, *product.Product) {
		return &product.NotFoundError{}, nil // Mock a not found product.
	}})

	// To get path params work is required to register them in a mux, instead of just call ServeHTTP on
	// handler.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}", rs.GetProduct(AlwaysValid[string]{}))

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf(
			"GetProduct should return a not found response: got status: %v msg: %v",
			rr.Code,
			rr.Body.String(),
		)
	}
}

func TestGetProduct_shouldReturnBadRequestIfInvalidID(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()

	rs := product.NewRoutes(&MockProductService{getFn: func(id string) (error, *product.Product) {
		return nil, &product.Product{}
	}})

	// Is necessary to use a ServeMux to get path values work.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}", rs.GetProduct(AlwaysInvalid[string]{}))

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf(
			"GetProduct should return a bad request response for an invalid id: got %v",
			rr.Code,
		)
	}
}
