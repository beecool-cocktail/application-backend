package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"gorm.io/gorm"
)

type cocktailMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLCocktailRepository(db *gorm.DB) domain.CocktailMySQLRepository {
	return &cocktailMySQLRepository{db}
}

func (c *cocktailMySQLRepository) GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination domain.PaginationMySQLRepository) ([]domain.Cocktail, int64, error) {
	var cocktail []domain.Cocktail
	var total int64

	orm := c.db.Model(&cocktail)

	for sort, dir := range pagination.SortByDir {
		dirString := sortbydir.ParseStringBySortByDir(dir)
		order := sortbydir.MakeSortAndDir(sort, dirString)
		orm.Order(order)
	}

	orm.Where(filter)
	orm.Count(&total)

	if pagination.PageSize != 0 {
		orm.Limit(pagination.PageSize).Offset((pagination.Page - 1)* pagination.PageSize)
	}

	res := orm.Find(&cocktail)

	return cocktail, total, res.Error
}

func (c *cocktailMySQLRepository) QueryByCocktailID(ctx context.Context, id int64) (domain.Cocktail, error) {
	var cocktail domain.Cocktail

	res := c.db.Where("cocktail_id = ?", id).Take(&cocktail)

	return cocktail, res.Error
}

func (c *cocktailMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, co *domain.Cocktail) error {

	res := tx.Select("cocktail_id", "user_id", "title", "description").Create(co)

	return res.Error
}
