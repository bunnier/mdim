package utils

import "fmt"

type AggregateError []error

func PrintAggregateError(errs AggregateError) {
	for _, err := range errs {
		fmt.Println(err)
	}
}
