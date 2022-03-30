package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type photoInfo struct {
	Photo        string `structs:"photo"`
	IsCoverPhoto bool   `structs:"is_cover_photo"`
}

type cocktailCocktailMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLCocktailStepRepository(db *gorm.DB) domain.CocktailPhotoMySQLRepository {
	return &cocktailCocktailMySQLRepository{db}
}

func (s *cocktailCocktailMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailPhoto) error {

	res := tx.Select("cocktail_id", "photo", "is_cover_photo").Create(c)

	return res.Error
}

func (s *cocktailCocktailMySQLRepository) QueryCoverPhotoByCocktailId(ctx context.Context, id int64) (string, error) {

	var photo domain.CocktailPhoto
	res := s.db.Select("photo").
		Where("cocktail_id = ?", id).
		Where("is_cover_photo = ?", true).
		Take(&photo)

	return photo.Photo, res.Error
}

func (s *cocktailCocktailMySQLRepository) QueryPhotosByCocktailId(ctx context.Context, id int64) ([]domain.CocktailPhoto, error) {

	var photos []domain.CocktailPhoto
	res := s.db.Select("id", "cocktail_id", "photo", "is_cover_photo", "created_date").
		Where("cocktail_id = ?", id).
		Find(&photos)

	return photos, res.Error
}

func (s *cocktailCocktailMySQLRepository) QueryPhotoById(ctx context.Context, id int64) (domain.CocktailPhoto, error) {

	var photo domain.CocktailPhoto
	res := s.db.Select("id", "cocktail_id", "photo", "is_cover_photo", "created_date").
		Where("id = ?", id).
		Take(&photo)

	return photo, res.Error
}

func (s *cocktailCocktailMySQLRepository) UpdateTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailPhoto) (int64, error) {
	var photo domain.CocktailPhoto
	updateColumn := photoInfo{
		Photo:        c.Photo,
		IsCoverPhoto: c.IsCoverPhoto,
	}

	res := tx.Model(&photo).Where("id = ?", c.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (s *cocktailCocktailMySQLRepository) DeleteByCocktailIDTx(ctx context.Context, tx *gorm.DB, id int64) error {
	var photo domain.CocktailPhoto

	res := tx.Where("cocktail_id = ?", id).Delete(&photo)

	return res.Error
}

func (s *cocktailCocktailMySQLRepository) DeleteByIDTx(ctx context.Context, tx *gorm.DB, id int64) error {
	var photo domain.CocktailPhoto

	res := tx.Where("id = ?", id).Delete(&photo)

	return res.Error
}
