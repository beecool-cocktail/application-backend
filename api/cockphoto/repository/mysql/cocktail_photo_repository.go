package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
)

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
