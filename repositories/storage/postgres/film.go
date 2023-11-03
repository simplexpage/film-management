package postgres

import (
	"context"
	"errors"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	"film-management/pkg/query"
	"film-management/pkg/query/filter"
	"film-management/pkg/query/pagination"
	"film-management/pkg/query/sort"
	"film-management/pkg/validation"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FilmRepository is a struct for Film.
type FilmRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewFilmRepository is a constructor for filmRepository.
func NewFilmRepository(db *gorm.DB, logger *zap.Logger) *FilmRepository {
	return &FilmRepository{
		db:     db,
		logger: logger,
	}
}

// CreateFilm is a method to create film.
func (f FilmRepository) CreateFilm(ctx context.Context, model *models.Film) error {
	if err := f.db.WithContext(ctx).Create(model).Error; err != nil {
		f.logger.Error("failed to add a new film in db", zap.Error(err))

		return domain.ErrFilmCreate
	}

	return nil
}

// UpdateFilm is a method to update film.
func (f FilmRepository) UpdateFilm(ctx context.Context, model *models.Film) error {
	if err := f.db.WithContext(ctx).Model(&model).Updates(model).Error; err != nil {
		f.logger.Error("failed to update a film in db", zap.Error(err))

		return domain.ErrFilmUpdate
	}

	return nil
}

// FindOneFilmByUUID is a method to find one film by UUID.
func (f FilmRepository) FindOneFilmByUUID(ctx context.Context, uuid uuid.UUID) (models.Film, error) {
	var film models.Film

	if result := f.db.WithContext(ctx).Where("uuid = ?", uuid).First(&film); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Film{}, validation.NotFoundError{Err: domain.ErrFilmNotFound}
		}

		f.logger.Error("failed to find ad in db", zap.Error(result.Error))

		return models.Film{}, domain.ErrFilmFind
	}

	return film, nil
}

// FindOneFilmByUUIDWithCreator is a method to find one film by UUID.
func (f FilmRepository) FindOneFilmByUUIDWithCreator(ctx context.Context, uuid uuid.UUID) (models.Film, error) {
	var film models.Film

	if result := f.db.WithContext(ctx).Preload("Creator").Where("uuid = ?", uuid).First(&film); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Film{}, validation.NotFoundError{Err: domain.ErrFilmNotFound}
		}

		f.logger.Error("failed to find ad in db", zap.Error(result.Error))

		return models.Film{}, domain.ErrFilmFind
	}

	return film, nil
}

// FindAllFilms is a method to find all films.
func (f FilmRepository) FindAllFilms(ctx context.Context, filterSortLimit query.FilterSortLimit) ([]models.Film, pagination.Pagination, error) {
	var films []models.Film

	condition := filter.GetDBFilterMapCondition(filterSortLimit.Filter)
	fmt.Println(condition)

	if result := f.db.WithContext(ctx).Where(condition).Limit(filterSortLimit.Limit).Offset(filterSortLimit.Offset).Order(sort.GetDBQueryForSort(filterSortLimit.Sort)).Find(&films); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, pagination.Pagination{}, nil
		}

		f.logger.Error("failed to find all films in db", zap.Error(result.Error))

		return nil, pagination.Pagination{}, domain.ErrFilmFindAll
	}

	var count int64
	if result := f.db.Model(models.Film{}).Where(condition).Count(&count); result.Error != nil {
		return nil, pagination.Pagination{}, result.Error
	}

	return films, pagination.NewPagination(int(count), filterSortLimit.Limit, filterSortLimit.Offset), nil
}

// DeleteFilm is a method to delete film.
func (f FilmRepository) DeleteFilm(ctx context.Context, uuid uuid.UUID) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Film{}).Error
	if err != nil {
		f.logger.Error("failed to delete film in db", zap.Error(err))

		return domain.ErrFilmDelete
	}

	return nil
}

// FilmExistsWithTitle checks if a film with the given title exists.
func (f FilmRepository) FilmExistsWithTitle(ctx context.Context, title string) error {
	var count int64

	// Check if a film with the same title exists
	err := f.db.WithContext(ctx).
		Model(&models.Film{}).
		Where("title = ?", title).
		Count(&count).
		Error

	if err != nil {
		f.logger.Error("failed to check film existence by title in db", zap.Error(err))

		return domain.ErrFilmCheckExistence
	}

	if count > 0 {
		return domain.ErrFilmExistsWithTitle
	}

	return nil
}
