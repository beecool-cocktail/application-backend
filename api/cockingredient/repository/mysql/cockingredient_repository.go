package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
)

type cocktailIngredientMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLCocktailIngredientRepository(db *gorm.DB) domain.CocktailIngredientMySQLRepository {
	return &cocktailIngredientMySQLRepository{db}
}

func (s *cocktailIngredientMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailIngredient) error {

	err := tx.Select("cocktail_id", "ingredient_name",
		"ingredient_amount", "ingredient_unit").Create(c).Error

	if err != nil {
		return err
	}

	return nil
}

