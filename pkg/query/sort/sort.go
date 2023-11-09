package sort

import (
	"errors"
	customError "film-management/pkg/errors"
	"fmt"
	"strings"
)

var (
	ErrMalformedSortQueryParameter = errors.New("malformed sort query parameter, should be field.orderdirection")
	ErrMalformedOrderDirection     = errors.New("malformed orderdirection in sort query parameter, should be asc or desc")
	ErrUnknownField                = errors.New("unknown field in sort query parameter")
)

const (
	OrderASC  = "asc"
	OrderDESC = "desc"
)

type Sortable interface {
	Field() string
	Order() string
}

type Opts struct {
	field string
	order string
}

// newOptions returns a Sortable object.
func newOptions(sortField string, availableSortFields []string) (*Opts, error) {
	splits := strings.Split(sortField, ".")

	if len(splits) != 2 {
		return nil, customError.ValidationError{Field: "sort", Err: ErrMalformedSortQueryParameter}
	}

	field, order := splits[0], splits[1]

	if order != OrderDESC && order != OrderASC {
		return nil, customError.ValidationError{Field: "sort", Err: ErrMalformedOrderDirection}
	}

	if !stringInSlice(availableSortFields, field) {
		return nil, customError.ValidationError{Field: "sort", Err: ErrUnknownField}
	}

	return &Opts{
		field: field,
		order: strings.ToUpper(order),
	}, nil
}

// Field returns the field name.
func (o *Opts) Field() string {
	return o.field
}

// Order returns the order direction in uppercase.
func (o *Opts) Order() string {
	return o.order
}

// stringInSlice checks if a string is in a slice of strings.
func stringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}

	return false
}

// GetSortOptions returns a Sortable object.
func GetSortOptions(sortField string, fields []string, sortByDefault string) (Sortable, error) {
	if sortField == "" {
		sortField = sortByDefault
	}

	return newOptions(sortField, fields)
}

// GetDBQueryForSort returns a string that can be used in a SQL query.
func GetDBQueryForSort(option Sortable) string {
	return fmt.Sprintf("%s %s", option.Field(), strings.ToUpper(option.Order()))
}
