package pagination

import (
	"errors"
	"film-management/pkg/validation"
	"math"
)

var (
	ErrLimitQueryParameterIsNoValidNumber  = errors.New("limit query parameter is no valid number")
	ErrOffsetQueryParameterIsNoValidNumber = errors.New("offset query parameter is no valid number")
)

// Pagination is a struct for pagination.
type Pagination struct {
	Page       int `json:"page" example:"1"`
	TotalPages int `json:"total_pages" example:"10"`
	PageSize   int `json:"page_size" example:"20"`
	TotalCount int `json:"total_count" example:"200"`
}

// NewPagination returns pagination.
func NewPagination(totalCount, limit, offset int) Pagination {
	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	if offset < 0 {
		offset = 0
	}

	page := offset/limit + 1

	if page > totalPages {
		page = totalPages
	}

	return Pagination{
		Page:       page,
		PageSize:   limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}

// GetLimitOption returns limit.
func GetLimitOption(limit int, defaultLimit int) (int, error) {
	if limit > 0 {
		if err := validateLimit(limit); err != nil {
			return 0, err
		}

		return limit, nil
	}

	return defaultLimit, nil
}

// GetOffsetOption returns offset.
func GetOffsetOption(offset int) (int, error) {
	if offset > 0 {
		if err := validateOffset(offset); err != nil {
			return 0, err
		}

		return offset, nil
	}

	return -1, nil
}

// validateOffset validates offset.
func validateOffset(offset int) error {
	if offset < -1 {
		return validation.CustomError{Field: "offset", Err: ErrOffsetQueryParameterIsNoValidNumber}
	}

	return nil
}

// validateLimit validates limit.
func validateLimit(limit int) error {
	if limit < -1 {
		return validation.CustomError{Field: "limit", Err: ErrLimitQueryParameterIsNoValidNumber}
	}

	return nil
}
