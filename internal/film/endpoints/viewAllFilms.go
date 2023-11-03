package endpoints

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"time"
)

// MakeViewAllFilmsEndpoint is an endpoint for ViewAllFilms.
func MakeViewAllFilmsEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(ViewAllFilmsRequest)
		if !ok {
			return ViewAllFilmsResponse{}, ErrInvalidRequest
		}

		if items, p, errViewAllFilms := s.ViewAllFilms(ctx, reqForm.FilterSortLimit); errViewAllFilms != nil {
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
	FilterSortLimit query.FilterSortLimit
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
	ReleaseDate string    `json:"release_date"`
	Cast        string    `json:"cast"`
	Genre       Genre     `json:"genre"`
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
			ReleaseDate: item.ReleaseDate.Format(time.DateOnly),
			Cast:        item.Cast,
			Genre:       GenreFromEnum(item.Genre),
			Synopsis:    item.Synopsis,
			CreatedAt:   time.Unix(item.CreatedAt, 0).Format(time.DateTime),
			UpdatedAt:   time.Unix(item.UpdatedAt, 0).Format(time.DateTime),
		})
	}

	return films
}
