package film

import (
	"context"
	"errors"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"film-management/pkg/query/sort"
	"film-management/pkg/validation"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Repository is a struct for work with film in db.
type Repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewRepository is a constructor for Repository.
func NewRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CreateFilm is a method to create film.
func (f Repository) CreateFilm(ctx context.Context, model *models.Film) error {
	if err := f.db.WithContext(ctx).Create(model).Error; err != nil {
		f.logger.Error("failed to add a new film in db", zap.Error(err))

		return domain.ErrFilmCreate
	}

	return nil
}

// UpdateFilm is a method to update film.
func (f Repository) UpdateFilm(ctx context.Context, model *models.Film) error {
	if err := f.db.WithContext(ctx).Model(&model).Updates(model).Error; err != nil {
		f.logger.Error("failed to update a film in db", zap.Error(err))

		return domain.ErrFilmUpdate
	}

	return nil
}

// FindOneFilmByUUID is a method to find one film by UUID.
func (f Repository) FindOneFilmByUUID(ctx context.Context, uuid uuid.UUID) (models.Film, error) {
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
func (f Repository) FindOneFilmByUUIDWithCreator(ctx context.Context, uuid uuid.UUID) (models.Film, error) {
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
func (f Repository) FindAllFilms(ctx context.Context, filterSortLimit query.FilterSortLimit) ([]models.Film, pagination.Pagination, error) {
	var films []models.Film

	// Build condition
	condition := f.db.Where("1 = 1")

	// Add filters to condition
	for field, value := range filterSortLimit.Filter {
		addFilmFiltersToCondition(condition, field, value, f)
	}

	// Find all films with condition
	if result := f.db.WithContext(ctx).
		Preload("Genres").
		Where(condition).
		Limit(filterSortLimit.Limit).
		Offset(filterSortLimit.Offset).
		Order(sort.GetDBQueryForSort(filterSortLimit.Sort)).
		Find(&films); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, pagination.Pagination{}, nil
		}

		f.logger.Error("failed to find all films in db", zap.Error(result.Error))

		return nil, pagination.Pagination{}, domain.ErrFilmFindAll
	}

	// Get count of films with condition for pagination
	var count int64

	if result := f.db.Model(models.Film{}).Where(condition).Count(&count); result.Error != nil {
		f.logger.Error("failed to get count of films in db", zap.Error(result.Error))

		return nil, pagination.Pagination{}, domain.ErrFilmGetCount
	}

	return films, pagination.NewPagination(int(count), filterSortLimit.Limit, filterSortLimit.Offset), nil
}

func addFilmFiltersToCondition(condition *gorm.DB, field string, value interface{}, f Repository) {
	switch field {
	case "title":
		condition = condition.Where("title LIKE ?", "%"+value.(string)+"%")
	case "release_date":
		// Check if value is []string or string
		if dates, ok := value.([]string); ok && len(dates) == 2 {
			condition = condition.Where("release_date BETWEEN ? AND ?", dates[0], dates[1])
		} else if date, ok := value.(string); ok {
			condition = condition.Where("release_date = ?", date)
		} else {
			f.logger.Debug("release_date is not []string or string")
		}
	case "genres":
		genreNames, ok := value.([]string)
		if !ok {
			f.logger.Debug("genres is not []string")
			return
		}
		// Get genre IDs
		var genreIDs []uint
		if result := f.db.Model(models.Genre{}).Where("LOWER(name) IN ?", genreNames).Pluck("id", &genreIDs); result.Error != nil {
			f.logger.Error("failed to get genre IDs", zap.Error(result.Error))
			return
		}
		// Add condition
		if len(genreIDs) > 0 {
			condition = condition.Where("EXISTS (SELECT 1 FROM film_genres WHERE films.uuid = film_genres.film_uuid AND film_genres.genre_id IN ?)", genreIDs)
		}
	default:
		f.logger.Debug("unknown filter field", zap.String("field", field))
	}
}

// DeleteFilm is a method to delete film.
func (f Repository) DeleteFilm(ctx context.Context, uuid uuid.UUID) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Film{}).Error
	if err != nil {
		f.logger.Error("failed to delete film in db", zap.Error(err))

		return domain.ErrFilmDelete
	}

	return nil
}

// FilmExistsWithTitle checks if a film with the given title exists.
func (f Repository) FilmExistsWithTitle(ctx context.Context, title string) error {
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

// CreateGenre создает новый жанр в базе данных.
func (f Repository) CreateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	if err := f.db.WithContext(ctx).Create(genre).Error; err != nil {
		f.logger.Error("failed to create genre", zap.Error(err))

		return nil, domain.ErrFilmCreateGenre
	}

	return genre, nil
}

// GetGenresByNames возвращает список жанров по их названиям.
func (f Repository) GetGenresByNames(ctx context.Context, names []string) ([]models.Genre, error) {
	var genres []models.Genre

	if err := f.db.WithContext(ctx).Where("name IN ?", names).Find(&genres).Error; err != nil {
		f.logger.Error("failed to get genres by names", zap.Error(err))

		return nil, domain.ErrFilmGetGenresByNames
	}

	return genres, nil
}
