package types

import (
	"fmt"
)

// AggregateError Aggregate a list of errors.
type AggregateError interface {
	// AddError Add a inner error.
	AddError(err error) AggregateError

	// Len Return the inner errors count.
	Len() int

	// PrintAggregateError Print all inner errors by line.
	PrintAggregateError()
}

// NewAggregateError Return a new AggregateError.
func NewAggregateError() AggregateError {
	return new(sliceAggregateError)
}

// A implement of AggregateError.
type sliceAggregateError struct {
	innerErrors []error
}

// AddError Add a inner error.
func (errs *sliceAggregateError) AddError(err error) AggregateError {
	if errs.innerErrors == nil {
		errs.innerErrors = make([]error, 0, 3)
	}
	errs.innerErrors = append(errs.innerErrors, err)
	return errs
}

// Len Return the inner errors count.
func (errs *sliceAggregateError) Len() int {
	return len(errs.innerErrors)
}

// PrintAggregateError Print all inner errors by line.
func (errs *sliceAggregateError) PrintAggregateError() {
	for _, err := range errs.innerErrors {
		fmt.Println(err)
	}
}
