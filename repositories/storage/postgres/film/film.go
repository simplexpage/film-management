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

func (f Repository) CreateFilm(ctx context.Context, model *models.Film) error {
	// Start a new transaction
	tx := f.db.Begin()

	// Check if the director with the specified name exists
	if err := f.createOrUpdateDirector(tx, model); err != nil {
		tx.Rollback()
		f.logger.Error("failed to update a film in db", zap.Error(err))

		return domain.ErrFilmUpdate
	}

	// Create the film
	if err := tx.Create(model).Error; err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	tx.Commit()

	return nil
}

func (f Repository) UpdateFilm(ctx context.Context, model *models.Film) error {
	tx := f.db.WithContext(ctx).Begin()

	// Check if the director with the specified name exists
	if err := f.createOrUpdateDirector(tx, model); err != nil {
		tx.Rollback()
		f.logger.Error("failed to update a film in db", zap.Error(err))

		return domain.ErrFilmUpdate
	}

	if err := tx.Model(&model).Updates(model).Error; err != nil {
		tx.Rollback()
		f.logger.Error("failed to update a film in db", zap.Error(err))

		return domain.ErrFilmUpdate
	}

	if err := tx.Model(&model).Association("Genres").Replace(model.Genres); err != nil {
		tx.Rollback()
		f.logger.Error("failed to update genres for the film in db", zap.Error(err))

		return domain.ErrFilmUpdate
	}

	if err := tx.Model(&model).Association("Casts").Replace(model.Casts); err != nil {
		tx.Rollback()
		f.logger.Error("failed to update casts for the film in db", zap.Error(err))

		return domain.ErrFilmUpdate
	}

	tx.Commit()

	return nil
}

// createOrUpdateDirector is a method to create or update director.
func (f Repository) createOrUpdateDirector(tx *gorm.DB, model *models.Film) error {
	var director models.Director
	err := tx.Where("name = ?", model.Director.Name).First(&director).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := tx.Create(&model.Director).Error
		if err != nil {
			return err
		}
		model.DirectorID = model.Director.ID
	} else {
		model.Director.ID = director.ID
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

// FindOneFilmForViewByUUID is a method to find one film by UUID.
func (f Repository) FindOneFilmForViewByUUID(ctx context.Context, uuid uuid.UUID) (models.Film, error) {
	var film models.Film

	if result := f.db.WithContext(ctx).
		Preload("Creator").
		Preload("Genres").
		Preload("Director").
		Preload("Casts").
		Where("uuid = ?", uuid).
		First(&film); result.Error != nil {
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
		Preload("Director").
		Preload("Casts").
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
			f.logger.Error("release_date is not []string or string")
		}
	case "genres":
		genreNames, ok := value.([]string)
		if !ok {
			f.logger.Error("genres is not []string")
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

// FilmExists checks if a film with the given filmID and title exists.
// The operation parameter specifies the type of operation: "add" or "update".
func (f Repository) FilmExists(ctx context.Context, title string, filmID uuid.UUID, operation models.Operation) error {
	var count int64

	switch operation {
	case models.OperationAdd:
		// Check if advertiser with the same userID exists
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
	case models.OperationUpdate:
		// Check if another film with the same title exists, except for the current film
		err := f.db.WithContext(ctx).
			Model(&models.Film{}).
			Where("title = ? AND uuid <> ?", title, filmID).
			Count(&count).
			Error

		if err != nil {
			f.logger.Error("failed to check film existence by title and uuid in db", zap.Error(err))

			return domain.ErrFilmCheckExistence
		}

		if count > 0 {
			return domain.ErrFilmExistsWithTitle
		}

	default:
		f.logger.Debug("unknown operation", zap.String("operation", string(operation)))

		return domain.ErrUnknownOperation
	}

	return nil
}

// CreateGenre creates a new genre.
func (f Repository) CreateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	if err := f.db.WithContext(ctx).Create(genre).Error; err != nil {
		f.logger.Error("failed to create genre", zap.Error(err))

		return nil, domain.ErrFilmCreateGenre
	}

	return genre, nil
}

// GetGenresByNames returns genres by names.
func (f Repository) GetGenresByNames(ctx context.Context, names []string) ([]models.Genre, error) {
	var genres []models.Genre

	if err := f.db.WithContext(ctx).Where("name IN ?", names).Find(&genres).Error; err != nil {
		f.logger.Error("failed to get genres by names", zap.Error(err))

		return nil, domain.ErrFilmGetGenresByNames
	}

	return genres, nil
}

// CreateCast creates a new cast.
func (f Repository) CreateCast(ctx context.Context, cast *models.Cast) (*models.Cast, error) {
	if err := f.db.WithContext(ctx).Create(cast).Error; err != nil {
		f.logger.Error("failed to create cast", zap.Error(err))

		return nil, domain.ErrFilmCreateCast
	}

	return cast, nil
}

// GetCastsByNames returns casts by names.
func (f Repository) GetCastsByNames(ctx context.Context, names []string) ([]models.Cast, error) {
	var casts []models.Cast

	if err := f.db.WithContext(ctx).Where("name IN ?", names).Find(&casts).Error; err != nil {
		f.logger.Error("failed to get casts by names", zap.Error(err))

		return nil, domain.ErrFilmGetCastsByNames
	}

	return casts, nil
}
