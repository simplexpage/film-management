package endpoints

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"film-management/pkg/query/sort"
	"film-management/pkg/validation"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"strings"
	"time"
)

// MakeViewAllFilmsEndpoint is an endpoint for ViewAllFilms.
func MakeViewAllFilmsEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(ViewAllFilmsRequest)
		if !ok {
			return ViewAllFilmsResponse{}, ErrInvalidRequest
		}

		// Validate form
		if errValidate := reqForm.Validate(); errValidate != nil {
			return ViewAllFilmsResponse{Err: errValidate}, nil
		}

		// Build FilterSortLimit
		builder := query.NewFilterSortLimitBuilder()

		// Get sort
		sortOption, err := sort.GetSortOptions(reqForm.Sort, []string{"title", "release_date"}, "release_date.desc")
		if err != nil {
			return ViewAllFilmsResponse{Err: err}, nil
		}

		// Get limit and offset
		limit, err := pagination.GetLimitOption(reqForm.Limit, 20)
		if err != nil {
			return ViewAllFilmsResponse{Err: err}, err
		}

		// Get offset from HTTP request
		offset, err := pagination.GetOffsetOption(reqForm.Offset)
		if err != nil {
			return ViewAllFilmsResponse{Err: err}, err
		}

		// Get filters from HTTP request
		myFilters, err := getFilterOptions(reqForm)
		if err != nil {
			return ViewAllFilmsResponse{Err: err}, nil
		}

		// Build FilterSortLimit
		filterSortLimit := builder.
			SetSort(sortOption).
			SetFilter(myFilters).
			SetLimit(limit).
			SetOffset(offset).
			Build()

		if items, p, errViewAllFilms := s.ViewAllFilms(ctx, filterSortLimit); errViewAllFilms != nil {
			return ViewAllFilmsResponse{Err: errViewAllFilms}, nil
		} else {
			return ViewAllFilmsResponse{
				Items:      domainAllFilmItemsToAllItemFilms(items),
				Pagination: p,
			}, nil
		}
	}
}

// ViewAllFilmsRequest is a request for ViewAllFilms.
type ViewAllFilmsRequest struct {
	Sort        string   `json:"sort" validate:"omitempty,min=3,max=30" example:"title.asc"`
	Limit       int      `json:"limit" validate:"omitempty,min=1,max=100" example:"10"`
	Offset      int      `json:"offset" validate:"omitempty,min=0" example:"0"`
	Title       string   `json:"title" validate:"omitempty,min=3,max=30" example:"Garry Potter"`
	ReleaseDate string   `json:"release_date" validate:"omitempty,customDate,customRangeDate" example:"2021-01-01,2021-12-31:2022-01-01"`
	Genres      []string `json:"genres" validate:"omitempty,min=1,max=5,dive,min=3,max=100" example:"action,adventure,sci-fi"`
}

// Validate is a method to validate form.
func (r *ViewAllFilmsRequest) Validate() error {
	// Get custom validator
	customValidator, err := validation.GetValidator()
	if err != nil {
		return err
	}

	// Register custom validator "customDate"
	err = customValidator.GetValidate().RegisterValidation("customDate", CustomDateValidator)
	if err != nil {
		return err
	}

	// Register custom validator "customRangeDate"
	err = customValidator.GetValidate().RegisterValidation("customRangeDate", CustomDateRangeValidator)
	if err != nil {
		return err
	}

	// Add translation for "customDate"
	err = customValidator.AddTranslation("customDate", fmt.Sprintf("{0} must be valid (YYYY-MM-DD or YYYY-MM-DD:YYYY-MM-DD)"))
	if err != nil {
		return err
	}

	// Add translation for "customRangeDate"
	err = customValidator.AddTranslation("customRangeDate", fmt.Sprintf("{0} must be valid the first date must be less than the second date"))
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// ViewAllFilmsResponse is a response for ViewAllFilms.
type ViewAllFilmsResponse struct {
	Items      []ItemAllFilms        `json:"items,omitempty"`
	Pagination pagination.Pagination `json:"pagination,omitempty"`
	Err        error                 `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r ViewAllFilmsResponse) Failed() error { return r.Err }

// ItemAllFilms is a response for ViewAllFilms.
type ItemAllFilms struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Genres      []string  `json:"genres"`
	ReleaseDate string    `json:"release_date"`
	Cast        string    `json:"cast"`
	Synopsis    string    `json:"synopsis"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// domainAllFilmItemsToAllItemFilms is a function to convert domain film items to all item films.
func domainAllFilmItemsToAllItemFilms(items []models.Film) []ItemAllFilms {
	films := make([]ItemAllFilms, 0, len(items))

	for _, item := range items {
		films = append(films, ItemAllFilms{
			UUID:        item.UUID,
			Title:       item.Title,
			Director:    item.Director,
			Genres:      convertGenresToStrings(item.Genres),
			ReleaseDate: item.ReleaseDate.Format(time.DateOnly),
			Cast:        item.Cast,
			Synopsis:    item.Synopsis,
			CreatedAt:   time.Unix(item.CreatedAt, 0).Format(time.DateTime),
			UpdatedAt:   time.Unix(item.UpdatedAt, 0).Format(time.DateTime),
		})
	}

	return films
}

// convertGenresToStrings is a function to convert genres to strings.
func convertGenresToStrings(genres []models.Genre) []string {
	genreNames := make([]string, len(genres))

	for i, genre := range genres {
		genreNames[i] = genre.Name
	}

	return genreNames
}

// getFilterOptions is a function to get filter options.
func getFilterOptions(reqForm ViewAllFilmsRequest) (query.Filter, error) {
	myFilter := make(query.Filter)

	// set release_date
	if reqForm.ReleaseDate != "" {
		if strings.Contains(reqForm.ReleaseDate, ":") {
			split := strings.Split(reqForm.ReleaseDate, ":")
			if len(split) != 2 {
				return nil, validation.CustomError{Field: "release_date", Err: fmt.Errorf("invalid date range format: %s", reqForm.ReleaseDate)}
			}
			myFilter["release_date"] = []string{split[0], split[1]}
		} else {
			myFilter["release_date"] = reqForm.ReleaseDate
		}
	}

	// set title
	if reqForm.Title != "" {
		myFilter["title"] = reqForm.Title
	}

	// set genres
	if len(reqForm.Genres) > 0 {
		myFilter["genres"] = reqForm.Genres
	}

	return myFilter, nil
}
