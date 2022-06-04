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

func (f *favoriteCocktailMySQLRepository) QueryCountsByUserID(ctx context.Context, id int64) (int64, error) {
	var cocktail domain.FavoriteCocktail
	var total int64

	orm := f.db.Model(&cocktail)
	orm.Where("user_id = ?", id)

	res := orm.Count(&total)

	return total, res.Error
}

func (f *favoriteCocktailMySQLRepository) DeleteTx(ctx context.Context, tx *gorm.DB, cocktailID, userID int64) error {
	var cocktail domain.FavoriteCocktail
	orm := tx.Model(&cocktail)
	if cocktailID > 0 {
		orm.Where("cocktail_id = ?", cocktailID)
	}

	if userID > 0 {
		orm.Where("user_id = ?", userID)
	}
	res := orm.Delete(&cocktail)

	return res.Error
}

func (f *favoriteCocktailMySQLRepository) Delete(ctx context.Context, cocktailID, userID int64) error {
	var cocktail domain.FavoriteCocktail

	res := f.db.Where("user_id = ? AND cocktail_id = ?", userID, cocktailID).Delete(&cocktail)

	return res.Error
}
