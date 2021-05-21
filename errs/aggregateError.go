package errs

import (
	"fmt"
)

// Aggregate a list of errors.
type AggregateError interface {
	// Add a inner error.
	AddError(err error) AggregateError

	// Return the inner errors count.
	Len() int

	// Print all inner errors by line.
	PrintAggregateError()
}

// Return a new AggregateError
func NewAggregateError() AggregateError {
	var _ AggregateError = new(sliceAggregateError) // force build check
	return new(sliceAggregateError)
}

// A implement of AggregateError.
type sliceAggregateError struct {
	innerErrors []error
}

// Add a inner error.
func (errs *sliceAggregateError) AddError(err error) AggregateError {
	if errs.innerErrors == nil {
		errs.innerErrors = make([]error, 0, 3)
	}
	errs.innerErrors = append(errs.innerErrors, err)
	return errs
}

// Return the inner errors count.
func (errs *sliceAggregateError) Len() int {
	return len(errs.innerErrors)
}

// Print all inner errors by line.
func (errs *sliceAggregateError) PrintAggregateError() {
	for _, err := range errs.innerErrors {
		fmt.Println(err)
	}
}
