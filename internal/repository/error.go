package repository

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"strings"
)

type dataError string

func (s dataError) Error() string {
	return string(s)
}

func (s dataError) Public() string {
	str := strings.Replace(string(s), "data: ", "", 1)
	return cases.Title(language.English).String(str)
}

type privateError string

func (s privateError) Error() string {
	return string(s)
}

var (
	// ErrNotFound is returned when a resource cannot be found
	ErrNotFound dataError = "data: resource not found"

	// ErrInvoiceNoRequired
	ErrInvoiceNoRequired dataError = "data: Invoice number is required"

	// ErrDateRequired
	ErrDateRequired dataError = "data: Date is required"

	// ErrEmailRequired
	ErrEmailRequired dataError = "data: Sender|Receiver email is required"

	// ErrNameRequired
	ErrNameRequired dataError = "data: Sender|Receiver name is required"

	// ErrCurrencyRequired
	ErrCurrencyRequired dataError = "data: currency is required"

	// ErrItemsRequired
	ErrItemsRequired dataError = "data: Invoice item(s) is required"

	// ErrIDInvalid is returned when an invalid ID is provided to a method like Delete.
	ErrIDInvalid privateError = "data: ID provided is invalid"

	// ErrServiceRequired is returned when a service is incorrect or not provided
	ErrServiceRequired privateError = "data: service is required"

	// ErrInvalidData  is returned when an invalid data is input
	ErrInvalidData dataError = "data: data provided is not valid"

)
