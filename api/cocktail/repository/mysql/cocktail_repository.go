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

func (c *cocktailMySQLRepository) GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination domain.PaginationMySQLRepository) (*[]domain.Cocktail, int64, error) {
	var cocktail []domain.Cocktail
	var total int64

	orm := c.db.Model(&cocktail)
	if pagination.PageSize != 0 {
		orm.Limit(pagination.PageSize).Offset((pagination.Page - 1)* pagination.PageSize)
	}

	for sort, dir := range pagination.SortByDir {
		dirString := sortbydir.ParseStringBySortByDir(dir)
		order := sortbydir.MakeSortAndDir(sort, dirString)
		orm.Order(order)
	}

	orm.Where(filter)
	res := orm.Find(&cocktail)
	orm.Count(&total)

	return &cocktail, total, res.Error
}
