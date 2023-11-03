package sort

import (
	"errors"
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

func NewOption(sortField string, availableSortFields []string) (*Opts, error) {
	splits := strings.Split(sortField, ".")

	if len(splits) != 2 {
		return nil, ErrMalformedSortQueryParameter
	}

	field, order := splits[0], splits[1]

	if order != OrderDESC && order != OrderASC {
		return nil, ErrMalformedOrderDirection
	}

	if !stringInSlice(availableSortFields, field) {
		return nil, ErrUnknownField
	}

	return &Opts{
		field: field,
		order: strings.ToUpper(order),
	}, nil
}

func (o *Opts) Field() string {
	return o.field
}
func (o *Opts) Order() string {
	return o.order
}

func stringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}

	return false
}

// GetDBQueryForSort returns a string that can be used in a SQL query
func GetDBQueryForSort(option Sortable) string {
	return fmt.Sprintf("%s %s", option.Field(), strings.ToUpper(option.Order()))
}
