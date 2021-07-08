package types

import (
	"fmt"
)

// AggregateError is a list of errors.
type AggregateError interface {
	// AddError add a inner error.
	AddError(err error) AggregateError

	// Len will return the inner errors count.
	Len() int

	// PrintAggregateError print all inner errors by line.
	PrintAggregateError()
}

// NewAggregateError return a new AggregateError.
func NewAggregateError() AggregateError {
	return new(sliceAggregateError)
}

// A implement of AggregateError.
type sliceAggregateError struct {
	innerErrors []error
}

// AddError add a inner error.
func (errs *sliceAggregateError) AddError(err error) AggregateError {
	if errs.innerErrors == nil {
		errs.innerErrors = make([]error, 0, 3)
	}
	errs.innerErrors = append(errs.innerErrors, err)
	return errs
}

// Len will return the inner errors count.
func (errs *sliceAggregateError) Len() int {
	return len(errs.innerErrors)
}

// PrintAggregateError print all inner errors by line.
func (errs *sliceAggregateError) PrintAggregateError() {
	for _, err := range errs.innerErrors {
		fmt.Println(err)
	}
}
