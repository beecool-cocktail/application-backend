package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"gorm.io/gorm"
)

type favoriteCocktailMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLFavoriteCocktailRepository(db *gorm.DB) domain.FavoriteCocktailMySQLRepository {
	return &favoriteCocktailMySQLRepository{db}
}

func (f *favoriteCocktailMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.FavoriteCocktail) error {

	err := tx.Select("cocktail_id", "user_id").Create(c).Error

	if err != nil {
		return err
	}

	return nil
}

func (f *favoriteCocktailMySQLRepository) QueryByUserID(ctx context.Context, id int64, pagination domain.PaginationMySQLRepository) ([]domain.FavoriteCocktail, int64, error) {
	var cocktails []domain.FavoriteCocktail
	var total int64

	orm := f.db.Model(&cocktails)

	for sort, dir := range pagination.SortByDir {
		dirString := sortbydir.ParseStringBySortByDir(dir)
		order := sortbydir.MakeSortAndDir(sort, dirString)
		orm.Order(order)
	}

	orm.Where("user_id = ?", id)
	orm.Count(&total)

	if pagination.PageSize != 0 {
		orm.Limit(pagination.PageSize).Offset((pagination.Page - 1) * pagination.PageSize)
	}

	res := orm.Find(&cocktails)

	return cocktails, total, res.Error
}

func (f *favoriteCocktailMySQLRepository) DeleteTx(ctx context.Context, tx *gorm.DB, cocktailID, userID int64) error {
	var cocktail domain.FavoriteCocktail

	res := tx.Where("user_id = ? AND cocktail_id = ?", userID, cocktailID).Delete(&cocktail)

	return res.Error
}