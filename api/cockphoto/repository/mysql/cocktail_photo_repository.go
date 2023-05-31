package mysql

import (
	"context"

	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type photoInfo struct {
	Photo        string `structs:"photo"`
	IsCoverPhoto bool   `structs:"is_cover_photo"`
	Order        int    `structs:"order"`
}

type photoOrder struct {
	IsCoverPhoto bool   `structs:"is_cover_photo"`
	Order int `structs:"order"`
}

type cocktailCocktailMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLCocktailStepRepository(db *gorm.DB) domain.CocktailPhotoMySQLRepository {
	return &cocktailCocktailMySQLRepository{db}
}

func (s *cocktailCocktailMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailPhoto) error {

	res := tx.Select("cocktail_id", "photo", "is_cover_photo", "low_quality_bundle_id", "is_low_quality", "order").Create(c)

	return res.Error
}

func (s *cocktailCocktailMySQLRepository) QueryCoverPhotoByCocktailId(ctx context.Context, id int64) (string, error) {

	var photo domain.CocktailPhoto
	res := s.db.Select("photo").
		Where("cocktail_id = ?", id).
		Where("is_cover_photo = ?", true).
		Where("is_low_quality = ?", false).
		Take(&photo)

	return photo.Photo, res.Error
}

func (s *cocktailCocktailMySQLRepository) QueryPhotosByCocktailId(ctx context.Context, id int64) ([]domain.CocktailPhoto, error) {

	var photos []domain.CocktailPhoto
	order := sortbydir.MakeSortAndDir("`order`", sortbydir.ASC.String())
	res := s.db.Select("id", "cocktail_id", "low_quality_bundle_id", "photo", "is_cover_photo", "created_date").
		Where("cocktail_id = ?", id).
		Where("is_low_quality = ?", false).
		Order(order).
		Find(&photos)

	return photos, res.Error
}

func (s *cocktailCocktailMySQLRepository) QueryLowQualityPhotosByCocktailId(ctx context.Context, id int64) ([]domain.CocktailPhoto, error) {

	var photos []domain.CocktailPhoto
	order := sortbydir.MakeSortAndDir("`order`", sortbydir.ASC.String())
	res := s.db.Select("id", "cocktail_id", "low_quality_bundle_id", "photo", "is_cover_photo", "created_date").
		Where("cocktail_id = ?", id).
		Where("is_low_quality = ?", true).
		Order(order).
		Find(&photos)

	return photos, res.Error
}

func (s *cocktailCocktailMySQLRepository) QueryPhotoById(ctx context.Context, id int64) (domain.CocktailPhoto, error) {

	var photo domain.CocktailPhoto
	res := s.db.Select("id", "cocktail_id", "low_quality_bundle_id", "photo", "is_cover_photo", "created_date").
		Where("id = ?", id).
		Take(&photo)

	return photo, res.Error
}

func (s *cocktailCocktailMySQLRepository) QueryLowQualityPhotoByBundleId(ctx context.Context, id int64) (domain.CocktailPhoto, error) {

	var photo domain.CocktailPhoto
	res := s.db.Select("id", "cocktail_id", "photo", "is_cover_photo", "created_date").
		Where("low_quality_bundle_id = ?", id).
		Where("is_low_quality = ?", true).
		Take(&photo)

	return photo, res.Error
}

func (s *cocktailCocktailMySQLRepository) UpdateTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailPhoto) (int64, error) {
	var photo domain.CocktailPhoto
	updateColumn := photoInfo{
		Photo:        c.Photo,
		Order:        c.Order,
		IsCoverPhoto: c.IsCoverPhoto,
	}

	res := tx.Model(&photo).Where("id = ?", c.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (s *cocktailCocktailMySQLRepository) UpdatePhotoOrderTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailPhoto) (int64, error) {
	var photo domain.CocktailPhoto
	updateColumn := photoOrder{
		Order: c.Order,
		IsCoverPhoto: c.IsCoverPhoto,
	}

	res := tx.Model(&photo).Where("low_quality_bundle_id = ?", c.LowQualityBundleID).Updates(structs.Map(updateColumn))

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

func (s *cocktailCocktailMySQLRepository) DeleteByLowQualityBundleIDTx(ctx context.Context, tx *gorm.DB, id int64) error {
	var photo domain.CocktailPhoto

	res := tx.Where("low_quality_bundle_id = ?", id).Delete(&photo)

	return res.Error
}
