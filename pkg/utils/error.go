package utils

import "strings"

type CombinedError struct {
	Errors []error
}

func (ce CombinedError) Error() string {
	errorMessages := make([]string, 0, len(ce.Errors))
	for _, err := range ce.Errors {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}
	return strings.Join(errorMessages, "\n")
}

func CombineErrors(errs ...error) error {
	nonNilErrors := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrors = append(nonNilErrors, err)
		}
	}
	if len(nonNilErrors) == 0 {
		return nil
	}
	return CombinedError{Errors: nonNilErrors}
}
