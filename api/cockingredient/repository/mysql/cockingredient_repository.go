package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type ingredientInfo struct {
	IngredientName   string `structs:"ingredient_name"`
	IngredientAmount string `structs:"ingredient_amount"`
}

type cocktailIngredientMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLCocktailIngredientRepository(db *gorm.DB) domain.CocktailIngredientMySQLRepository {
	return &cocktailIngredientMySQLRepository{db}
}

func (s *cocktailIngredientMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailIngredient) error {

	err := tx.Select("cocktail_id", "ingredient_name",
		"ingredient_amount").Create(c).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *cocktailIngredientMySQLRepository) QueryByCocktailId(ctx context.Context, id int64) ([]domain.CocktailIngredient, error) {

	var ingredients []domain.CocktailIngredient
	err := s.db.Select("id", "ingredient_name", "ingredient_amount").
		Where("cocktail_id = ?", id).
		Find(&ingredients).Error

	if err != nil {
		return []domain.CocktailIngredient{}, err
	}

	return ingredients, nil
}

func (s *cocktailIngredientMySQLRepository) UpdateTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailIngredient) (int64, error) {
	var ingredient domain.CocktailIngredient
	updateColumn := ingredientInfo{
		IngredientName:   c.IngredientName,
		IngredientAmount: c.IngredientAmount,
	}

	res := tx.Model(&ingredient).Where("id = ?", c.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (s *cocktailIngredientMySQLRepository) DeleteByCocktailIDTx(ctx context.Context, tx *gorm.DB, id int64) error {
	var ingredient domain.CocktailIngredient

	res := tx.Where("cocktail_id = ?", id).Delete(&ingredient)

	return res.Error
}
