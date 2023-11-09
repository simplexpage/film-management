package film

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	customError "film-management/pkg/errors"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"film-management/pkg/query/sort"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Repository is a struct for work with film in db.
type Repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewFilmRepository is a constructor for Repository.
func NewFilmRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (f Repository) CreateFilm(ctx context.Context, model *models.Film) error {
	// Start a new transaction
	tx := f.db.WithContext(ctx).Begin()

	// Check if the director with the specified name exists and create or update it
	if err := f.createOrUpdateDirector(tx, model); err != nil {
		tx.Rollback()
		f.logger.Error("filmRepo.CreateFilm.createOrUpdateDirector", zap.Error(err))

		return errors.Wrap(err, "filmRepo.CreateFilm.createOrUpdateDirector")
	}

	// Create the film
	if err := tx.Create(model).Error; err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		f.logger.Error("filmRepo.CreateFilm.Create", zap.Error(err))

		return errors.Wrap(err, "filmRepo.CreateFilm.Create")
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
		f.logger.Error("filmRepo.UpdateFilm.createOrUpdateDirector", zap.Error(err))

		return errors.Wrap(err, "filmRepo.UpdateFilm.createOrUpdateDirector")
	}

	// Update the film
	if err := tx.Model(&model).Updates(model).Error; err != nil {
		tx.Rollback()
		f.logger.Error("filmRepo.UpdateFilm.Updates", zap.Error(err))

		return errors.Wrap(err, "filmRepo.UpdateFilm.Updates")
	}

	// Replace genres and casts
	if err := tx.Model(&model).Association("Genres").Replace(model.Genres); err != nil {
		tx.Rollback()
		f.logger.Error("filmRepo.UpdateFilm.ReplaceGenres", zap.Error(err))

		return errors.Wrap(err, "filmRepo.UpdateFilm.ReplaceGenres")
	}

	if err := tx.Model(&model).Association("Casts").Replace(model.Casts); err != nil {
		tx.Rollback()
		f.logger.Error("filmRepo.UpdateFilm.ReplaceCasts", zap.Error(err))

		return errors.Wrap(err, "filmRepo.UpdateFilm.ReplaceCasts")
	}

	tx.Commit()

	return nil
}

// createOrUpdateDirector is a method to create or update director.
func (f Repository) createOrUpdateDirector(tx *gorm.DB, model *models.Film) error {
	// Check if the director with the specified name exists
	var director models.Director
	err := tx.Where("name = ?", model.Director.Name).First(&director).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "filmRepo.createOrUpdateDirector.First")
	}
	// If the director with the specified name does not exist, then create it
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errCreate := tx.Create(&model.Director).Error
		if errCreate != nil {
			f.logger.Error("filmRepo.createOrUpdateDirector.Create", zap.Error(errCreate))

			return errors.Wrap(errCreate, "filmRepo.createOrUpdateDirector.Create")
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
			return models.Film{}, errors.Wrap(domain.ErrFilmNotFound, "filmRepo.FindOneFilmByUUID.First")
		}

		f.logger.Error("filmRepo.FindOneFilmByUUID.First", zap.Error(result.Error))

		return models.Film{}, errors.Wrap(result.Error, "filmRepo.FindOneFilmByUUID.First")
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
			return models.Film{}, errors.Wrap(domain.ErrFilmNotFound, "filmRepo.FindOneFilmForViewByUUID.First")
		}

		f.logger.Error("filmRepo.FindOneFilmForViewByUUID.First", zap.Error(result.Error))

		return models.Film{}, errors.Wrap(result.Error, "filmRepo.FindOneFilmForViewByUUID.First")
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
		if err := addFilmFiltersToCondition(condition, field, value, f); err != nil {
			f.logger.Error("filmRepo.FindAllFilms.addFilmFiltersToCondition", zap.Error(err))

			return nil, pagination.Pagination{}, err
		}
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

		f.logger.Error("filmRepo.FindAllFilms.Find", zap.Error(result.Error))

		return nil, pagination.Pagination{}, domain.ErrFilmFindAll
	}

	// Get count of films with condition for pagination
	var count int64

	if result := f.db.Model(models.Film{}).Where(condition).Count(&count); result.Error != nil {
		f.logger.Error("filmRepo.FindAllFilms.Count", zap.Error(result.Error))

		return nil, pagination.Pagination{}, domain.ErrFilmFindAll
	}

	return films, pagination.NewPagination(int(count), filterSortLimit.Limit, filterSortLimit.Offset), nil
}

// addFilmFiltersToCondition is a method to add film filters to condition.
func addFilmFiltersToCondition(condition *gorm.DB, field string, value interface{}, f Repository) error {
	switch field {
	case "title":
		return addTitleFilter(condition, value)
	case "release_date":
		return addReleaseDateFilter(condition, value)
	case "genres":
		return addGenresFilter(condition, value, f)
	default:
		return customError.ValidationError{Field: field, Err: domain.ErrFilmUnknownField}
	}
}

// addTitleFilter is a method to add title filter.
func addTitleFilter(condition *gorm.DB, value interface{}) error {
	title, ok := value.(string)
	if !ok {
		return customError.ValidationError{Field: "title", Err: domain.ErrFilmFilterWrong}
	}
	condition = condition.Where("title LIKE ?", "%"+title+"%")

	return nil
}

// addReleaseDateFilter is a method to add release date filter.
func addReleaseDateFilter(condition *gorm.DB, value interface{}) error {
	// Check if value is []string or string
	dates, ok := value.([]string)
	if !ok || len(dates) != 2 {
		return customError.ValidationError{Field: "release_date", Err: domain.ErrFilmFilterWrong}
	}
	condition = condition.Where("release_date BETWEEN ? AND ?", dates[0], dates[1])

	return nil
}

// addGenresFilter is a method to add genres filter.
func addGenresFilter(condition *gorm.DB, value interface{}, f Repository) error {
	genreNames, ok := value.([]string)
	if !ok {
		return customError.ValidationError{Field: "genres", Err: domain.ErrFilmFilterWrong}
	}

	// Get genre IDs
	var genreIDs []uint
	if result := f.db.Model(models.Genre{}).Where("LOWER(name) IN ?", genreNames).Pluck("id", &genreIDs); result.Error != nil {
		f.logger.Error("filmRepo.addFilmFiltersToCondition.Pluck", zap.Error(result.Error))

		return domain.ErrFilmFindGenres
	}

	// Add condition
	if len(genreIDs) > 0 {
		condition = condition.Where("EXISTS (SELECT 1 FROM film_genres WHERE films.uuid = film_genres.film_uuid AND film_genres.genre_id IN ?)", genreIDs)
	} else {
		return customError.ValidationError{Field: "genres", Err: domain.ErrFilmGenresNotFound}
	}

	return nil
}

// DeleteFilm is a method to delete film.
func (f Repository) DeleteFilm(ctx context.Context, uuid uuid.UUID) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Film{}).Error
	if err != nil {
		f.logger.Error("filmRepo.DeleteFilm.Delete", zap.Error(err))

		return errors.Wrap(err, "filmRepo.DeleteFilm.Delete")
	}

	return nil
}

// FilmExistsWithTitle checks if a film with the given filmID and title exists.
// The operation parameter specifies the type of operation: "add" or "update".
func (f Repository) FilmExistsWithTitle(ctx context.Context, title string, filmID uuid.UUID, operation models.Operation) error {
	var count int64

	switch operation {
	case models.OperationAdd:
		// Check if a film with the same title exists
		err := f.db.WithContext(ctx).
			Model(&models.Film{}).
			Where("title = ?", title).
			Count(&count).
			Error

		if err != nil {
			f.logger.Error("filmRepo.FilmExists.OperationAdd.Count", zap.Error(err))

			return errors.Wrap(err, "filmRepo.FilmExists.OperationAdd.Count")
		}

		// If count > 0, then a film with the same title exists
		if count > 0 {
			return errors.Wrap(domain.ErrFilmExistsWithTitle, "filmRepo.FilmExists.OperationAdd.Count")
		}
	case models.OperationUpdate:
		// Check if another film with the same title exists, except for the current film
		err := f.db.WithContext(ctx).
			Model(&models.Film{}).
			Where("title = ? AND uuid <> ?", title, filmID).
			Count(&count).
			Error

		if err != nil {
			f.logger.Error("filmRepo.FilmExists.OperationUpdate.Count", zap.Error(err))

			return errors.Wrap(err, "filmRepo.FilmExists.OperationUpdate.Count")
		}

		// If count > 0, then a film with the same title exists
		if count > 0 {
			return errors.Wrap(domain.ErrFilmExistsWithTitle, "filmRepo.FilmExists.OperationUpdate.Count")
		}

	default:
		f.logger.Error("filmRepo.FilmExists.unknown operation", zap.String("operation", string(operation)))

		return errors.Wrap(domain.ErrFilmExistsWithTitle, "filmRepo.FilmExists.unknown operation")
	}

	return nil
}

// CreateGenre creates a new genre.
func (f Repository) CreateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	if err := f.db.WithContext(ctx).Create(genre).Error; err != nil {
		f.logger.Error("filmRepo.CreateGenre.Create", zap.Error(err))

		return nil, errors.Wrap(err, "filmRepo.CreateGenre.Create")
	}

	return genre, nil
}

// GetGenresByNames returns genres by names.
func (f Repository) GetGenresByNames(ctx context.Context, names []string) ([]models.Genre, error) {
	var genres []models.Genre

	if err := f.db.WithContext(ctx).Where("name IN ?", names).Find(&genres).Error; err != nil {
		f.logger.Error("filmRepo.GetGenresByNames.Find", zap.Error(err))

		return nil, errors.Wrap(err, "filmRepo.GetGenresByNames.Find")
	}

	return genres, nil
}

// CreateCast creates a new cast.
func (f Repository) CreateCast(ctx context.Context, cast *models.Cast) (*models.Cast, error) {
	if err := f.db.WithContext(ctx).Create(cast).Error; err != nil {
		f.logger.Error("filmRepo.CreateCast.Create", zap.Error(err))

		return nil, errors.Wrap(err, "filmRepo.CreateCast.Create")
	}

	return cast, nil
}

// GetCastsByNames returns casts by names.
func (f Repository) GetCastsByNames(ctx context.Context, names []string) ([]models.Cast, error) {
	var casts []models.Cast

	if err := f.db.WithContext(ctx).Where("name IN ?", names).Find(&casts).Error; err != nil {
		f.logger.Error("filmRepo.GetCastsByNames.Find", zap.Error(err))

		return nil, errors.Wrap(err, "filmRepo.GetCastsByNames.Find")
	}

	return casts, nil
}
