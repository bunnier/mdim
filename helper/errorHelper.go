package helper

import "fmt"

func PrintAggregateError(errs []error) {
	for _, err := range errs {
		fmt.Println(err)
	}
}
