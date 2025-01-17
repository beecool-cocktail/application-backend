package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/cockarticletype"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type cocktailInfo struct {
	Title       string `structs:"title"`
	Description string `structs:"description"`
}

type cocktailCategory struct {
	Category int `structs:"category"`
}

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
		orm.Limit(pagination.PageSize).Offset((pagination.Page - 1) * pagination.PageSize)
	}

	res := orm.Find(&cocktail)

	return cocktail, total, res.Error
}

func (c *cocktailMySQLRepository) QueryByCocktailID(ctx context.Context, id int64) (domain.Cocktail, error) {
	var cocktail domain.Cocktail

	res := c.db.Where("cocktail_id = ?", id).Take(&cocktail)

	return cocktail, res.Error
}

func (c *cocktailMySQLRepository) QueryFormalByUserID(ctx context.Context, id int64) ([]domain.Cocktail, error) {
	var cocktail []domain.Cocktail

	res := c.db.Where("user_id = ? AND category = ?", id, cockarticletype.Formal).Find(&cocktail)

	return cocktail, res.Error
}

func (c *cocktailMySQLRepository) QueryFormalCountsByUserID(ctx context.Context, id int64) (int64, error) {
	var cocktail []domain.Cocktail
	var total int64

	orm := c.db.Model(&cocktail)
	orm.Where("user_id = ?", id)
	orm.Where("category = ?", cockarticletype.Formal)

	res := orm.Count(&total)

	return total, res.Error
}

func (c *cocktailMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, co *domain.Cocktail) error {

	res := tx.Select("cocktail_id", "user_id", "title", "description", "category").Create(co)

	return res.Error
}

func (c *cocktailMySQLRepository) UpdateTx(ctx context.Context, tx *gorm.DB, co *domain.Cocktail) (int64, error) {
	var cocktail domain.Cocktail
	updateColumn := cocktailInfo{
		Title:       co.Title,
		Description: co.Description,
	}

	res := tx.Model(&cocktail).Where("cocktail_id = ?", co.CocktailID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (c *cocktailMySQLRepository) UpdateCategoryTx(ctx context.Context, tx *gorm.DB, co *domain.Cocktail) (int64, error) {
	var cocktail domain.Cocktail
	updateColumn := cocktailCategory{
		Category: co.Category,
	}

	res := tx.Model(&cocktail).Where("cocktail_id = ?", co.CocktailID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (c *cocktailMySQLRepository) IncreaseNumberOfCollectionTx(ctx context.Context, tx *gorm.DB,
	cocktailID int64) (int64, error) {
	var cocktail domain.Cocktail

	res := tx.Model(&cocktail).Where("cocktail_id = ?", cocktailID).
		UpdateColumn("number_of_collection", gorm.Expr("number_of_collection + 1"))

	return res.RowsAffected, res.Error
}

func (c *cocktailMySQLRepository) DecreaseNumberOfCollectionTx(ctx context.Context, tx *gorm.DB,
	cocktailID int64) (int64, error) {
	var cocktail domain.Cocktail

	res := tx.Model(&cocktail).Where("cocktail_id = ?", cocktailID).
		UpdateColumn("number_of_collection", gorm.Expr("number_of_collection - 1"))

	return res.RowsAffected, res.Error
}

func (c *cocktailMySQLRepository) DeleteTx(ctx context.Context, tx *gorm.DB, id int64) error {
	var cocktail domain.Cocktail

	res := tx.Where("cocktail_id = ?", id).Delete(&cocktail)

	return res.Error
}
