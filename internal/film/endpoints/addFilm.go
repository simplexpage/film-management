package endpoints

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	"film-management/pkg/validation"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"strings"
	"time"
)

// MakeAddFilmEndpoint is an endpoint for AddAd.
func MakeAddFilmEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(AddFilmRequest)
		if !ok {
			return AddFilmResponse{}, ErrInvalidRequest
		}

		// Validate form
		if errValidate := reqForm.Validate(); errValidate != nil {
			return AddFilmResponse{Err: errValidate}, nil
		}

		// Parse date
		parseDate, err := time.Parse(time.DateOnly, reqForm.ReleaseDate)
		if err != nil {
			return AddFilmResponse{Err: err}, nil
		}

		// Parse creator UUID
		parseCreatorUUID, err := uuid.Parse(reqForm.CreatorID)
		if err != nil {
			return AddFilmResponse{Err: err}, nil
		}

		// Prepare genres
		var genres []models.Genre
		for _, genreName := range reqForm.Genres {
			genre := models.Genre{Name: strings.ToLower(genreName)}
			genres = append(genres, genre)
		}

		// Prepare casts
		var casts []models.Cast
		for _, castName := range reqForm.Casts {
			cast := models.Cast{Name: castName}
			casts = append(casts, cast)
		}

		// Prepare a film model
		model := &models.Film{
			CreatorID:   parseCreatorUUID,
			Title:       reqForm.Title,
			Director:    models.Director{Name: reqForm.Director},
			ReleaseDate: parseDate,
			Casts:       casts,
			Synopsis:    reqForm.Synopsis,
			Genres:      genres,
		}

		// Add film
		if errAddFilm := s.AddFilm(ctx, model); errAddFilm != nil {
			return AddFilmResponse{Err: errAddFilm}, nil
		}

		return AddFilmResponse{
			Item: domainFilmToItemFilm(model),
		}, nil
	}
}

// AddFilmRequest is a request for Add film.
type AddFilmRequest struct {
	CreatorID string `json:"creatorID" validate:"required,uuid4" swaggerignore:"true"`

	Title       string   `json:"title" validate:"required,min=3,max=100" example:"Garry Potter"`
	Director    string   `json:"director" validate:"required,min=3,max=40" example:"John Doe"`
	ReleaseDate string   `json:"releaseDate" validate:"required,customDate" example:"2021-01-01"`
	Genres      []string `json:"genres" validate:"required,min=1,max=5,dive,min=3,max=100" example:"action,adventure,sci-fi"`
	Casts       []string `json:"casts" validate:"required,min=1,max=10,dive,min=3,max=100" example:"John Doe, Jane Doe, Foo Bar, Baz Quux"`
	Synopsis    string   `json:"synopsis" validate:"required,min=10,max=1000" example:"This is a synopsis."`
}

// Validate is a method to validate form.
func (r *AddFilmRequest) Validate() error {
	// Get custom validator
	customValidator, err := validation.GetValidator()
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// AddFilmResponse is a response for AddFilm.
type AddFilmResponse struct {
	Item ItemFilm `json:"item,omitempty"`
	Err  error    `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r AddFilmResponse) Failed() error { return r.Err }

// ItemFilm is a response for ViewFilm.
type ItemFilm struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Genres      []string  `json:"genres"`
	ReleaseDate string    `json:"release_date"`
	Casts       []string  `json:"casts"`
	Synopsis    string    `json:"synopsis"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// domainFilmToItemFilm is a method to convert domain Film to Item Film.
func domainFilmToItemFilm(item *models.Film) ItemFilm {
	return ItemFilm{
		UUID:        item.UUID,
		Title:       item.Title,
		Director:    item.Director.Name,
		Genres:      convertGenresToStrings(item.Genres),
		ReleaseDate: item.ReleaseDate.Format(time.DateOnly),
		Casts:       convertCastsToStrings(item.Casts),
		Synopsis:    item.Synopsis,
		CreatedAt:   time.Unix(item.CreatedAt, 0).Format(time.DateTime),
		UpdatedAt:   time.Unix(item.UpdatedAt, 0).Format(time.DateTime),
	}
}
