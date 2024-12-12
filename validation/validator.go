package validation

// Validator abstracts a validator for type T.
type Validator[T any] interface {
	// Validate check if a value `v` is valid.
	Validate(v *T) error
}
