package endpoints

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	customError "film-management/pkg/errors"
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
			return ViewAllFilmsResponse{}, customError.ErrInvalidRequest
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
	ReleaseDate string   `json:"release_date" validate:"omitempty,customRangeDate,customRangeDateCorrect" example:"2021-01-01,2021-12-31:2022-01-01"`
	Genres      []string `json:"genres" validate:"omitempty,min=1,max=5,dive,min=3,max=100" example:"action,adventure,sci-fi"`
}

// Validate is a method to validate form.
func (r *ViewAllFilmsRequest) Validate() error {
	// Get custom validator
	customValidator, err := validation.GetValidator()
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// ViewAllFilmsResponse is a response for ViewAllFilms.
type ViewAllFilmsResponse struct {
	Items      []ItemAllFilms        `json:"items"`
	Pagination pagination.Pagination `json:"pagination,omitempty"`
	Err        error                 `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r ViewAllFilmsResponse) Failed() error { return r.Err }

// ItemAllFilms is a response for ViewAllFilms.
type ItemAllFilms struct {
	UUID        uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string    `json:"title" example:"Garry Potter"`
	Director    string    `json:"director" example:"John Doe"`
	Genres      []string  `json:"genres" example:"action,adventure,sci-fi"`
	ReleaseDate string    `json:"release_date" example:"2021-01-01"`
	Casts       []string  `json:"casts" example:"John Doe,Jane Doe,Foo Bar"`
	Synopsis    string    `json:"synopsis" example:"This is a synopsis."`
	CreatedAt   string    `json:"created_at" example:"2021-01-01 00:00:00"`
	UpdatedAt   string    `json:"updated_at" example:"2021-01-01 00:00:00"`
}

// domainAllFilmItemsToAllItemFilms is a function to convert domain film items to all item films.
func domainAllFilmItemsToAllItemFilms(items []models.Film) []ItemAllFilms {
	films := make([]ItemAllFilms, 0, len(items))

	for _, item := range items {
		films = append(films, ItemAllFilms{
			UUID:        item.UUID,
			Title:       item.Title,
			Director:    item.Director.Name,
			Genres:      convertGenresToStrings(item.Genres),
			ReleaseDate: item.ReleaseDate.Format(time.DateOnly),
			Casts:       convertCastsToStrings(item.Casts),
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

// convertCastsToStrings is a function to convert casts to strings.
func convertCastsToStrings(casts []models.Cast) []string {
	castNames := make([]string, len(casts))

	for i, cast := range casts {
		castNames[i] = cast.Name
	}

	return castNames
}

// getFilterOptions is a function to get filter options.
func getFilterOptions(reqForm ViewAllFilmsRequest) (query.Filter, error) {
	myFilter := make(query.Filter)

	// set release_date
	if reqForm.ReleaseDate != "" {
		if strings.Contains(reqForm.ReleaseDate, ":") {
			split := strings.Split(reqForm.ReleaseDate, ":")
			if len(split) != 2 {
				return nil, customError.ValidationError{Field: "release_date", Err: fmt.Errorf("invalid date range format: %s", reqForm.ReleaseDate)}
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
