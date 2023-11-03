package query

import (
	"errors"
	"film-management/pkg/query/filter"
	"film-management/pkg/query/sort"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrLimitQueryParameterIsNoValidNumber  = errors.New("limit query parameter is no valid number")
	ErrOffsetQueryParameterIsNoValidNumber = errors.New("offset query parameter is no valid number")
)

// FilterSortLimit is a struct for FilterSortLimit.
type FilterSortLimit struct {
	Sort   sort.Sortable
	Filter []filter.Filter
	Limit  int
	Offset int
}

// NewFilterSortLimitFromHTTPRequest returns FilterSortLimit from HTTP request.
func NewFilterSortLimitFromHTTPRequest(r *http.Request, model interface{}, sortByDefault string) (FilterSortLimit, error) {
	// Get sort option
	sortOption, err := getSortOption(r, model, sortByDefault)
	if err != nil {
		return FilterSortLimit{}, err
	}

	// Get limit
	limit, err := getLimit(r)
	if err != nil {
		return FilterSortLimit{}, err
	}

	// Get offset
	offset, err := getOffset(r)
	if err != nil {
		return FilterSortLimit{}, err
	}

	// Get filters
	filters, err := getFiltersFromHTTPRequest(r, model)
	if err != nil {
		return FilterSortLimit{}, err
	}

	return FilterSortLimit{Limit: limit, Offset: offset, Sort: sortOption, Filter: filters}, nil
}

// getSortOption returns sort option.
func getSortOption(r *http.Request, model interface{}, sortByDefault string) (sort.Sortable, error) {
	availableFieldsForSortAndFilterFromModel := getFieldsFromModel(model)

	sortField := r.URL.Query().Get("sort")
	if sortField == "" {
		sortField = sortByDefault
	}

	sortOption, err := sort.NewOption(sortField, availableFieldsForSortAndFilterFromModel)
	return sortOption, err
}

// getLimit returns limit.
func getLimit(r *http.Request) (int, error) {
	limit := 20
	limitField := r.URL.Query().Get("limit")
	if limitField != "" {
		limit, err := strconv.Atoi(limitField)
		if err != nil {
			return 0, err
		}
		if err := validateLimit(limit); err != nil {
			return 0, err
		}
		return limit, nil
	}
	return limit, nil
}

// getOffset returns offset.
func getOffset(r *http.Request) (int, error) {
	offset := -1
	offsetField := r.URL.Query().Get("offset")
	if offsetField != "" {
		offset, err := strconv.Atoi(offsetField)
		if err != nil {
			return 0, err
		}
		if err := validateOffset(offset); err != nil {
			return 0, err
		}
		return offset, nil
	}
	return offset, nil
}

// getFiltersFromHTTPRequest returns filters from HTTP request.
func getFiltersFromHTTPRequest(r *http.Request, model interface{}) ([]filter.Filter, error) {
	var fields []filter.Filter
	v := reflect.ValueOf(model)

	for i := 0; i < v.Type().NumField(); i++ {
		field := v.Type().Field(i).Tag.Get("json")
		fieldType := v.Type().Field(i).Type.String()
		filterField := r.URL.Query().Get(field)

		if filterField == "" {
			continue
		}

		f, err := createFilter(field, fieldType, filterField)
		if err != nil {
			return nil, err
		}

		fields = append(fields, f)
	}

	return fields, nil
}

// createFilter returns filter.
func createFilter(field, fieldType, filterField string) (filter.Filter, error) {
	switch fieldType {
	case "string", "uuid.UUID":
		return filter.Filter{Field: field, Condition: "LIKE", Value: filterField}, nil
	case "int", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "uint", "float32", "float64":
		return filter.Filter{Field: field, Condition: "=", Value: filterField, ValueInt: true}, nil
	case "time.Time":
		return createTimeFilter(field, filterField)
	default:
		return filter.Filter{}, fmt.Errorf("undefined field type %s", fieldType)
	}
}

// createTimeFilter returns time filter.
func createTimeFilter(field, filterField string) (filter.Filter, error) {
	if strings.Contains(filterField, ":") {
		split := strings.Split(filterField, ":")
		if len(split) != 2 {
			return filter.Filter{}, fmt.Errorf("invalid date range format: %s", filterField)
		}
		// Parse from date
		parseFromDate, err := time.Parse(time.DateOnly, split[0])
		if err != nil {
			return filter.Filter{}, err
		}

		// Parse to date
		parseToDate, err := time.Parse(time.DateOnly, split[1])
		if err != nil {
			return filter.Filter{}, err
		}

		return filter.Filter{
			Field:     field,
			Condition: "BETWEEN",
			Value:     parseFromDate.Format(time.DateOnly),
			Value2:    parseToDate.Format(time.DateOnly),
		}, nil
	} else {
		// Parse date
		parseDate, err := time.Parse(time.DateOnly, filterField)
		if err != nil {
			return filter.Filter{}, err
		}

		return filter.Filter{Field: field, Condition: "=", Value: parseDate.Format(time.DateOnly)}, nil
	}
}

// getFieldsFromModel returns fields from model.
func getFieldsFromModel(model interface{}) []string {
	var fields []string

	v := reflect.ValueOf(model)

	for i := 0; i < v.Type().NumField(); i++ {
		fields = append(fields, v.Type().Field(i).Tag.Get("json"))
	}

	return fields
}

// validateOffset validates offset.
func validateOffset(offset int) error {
	if offset < -1 {
		return ErrOffsetQueryParameterIsNoValidNumber
	}

	return nil
}

// validateLimit validates limit.
func validateLimit(limit int) error {
	if limit < -1 {
		return ErrLimitQueryParameterIsNoValidNumber
	}

	return nil
}
