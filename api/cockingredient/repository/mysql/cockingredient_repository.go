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

func (s *cocktailIngredientMySQLRepository) QueryByCocktailId(ctx context.Context, id int64) ([]domain.CocktailIngredient, error) {

	var ingredients []domain.CocktailIngredient
	err := s.db.Select("id", "ingredient_name", "ingredient_amount", "ingredient_unit").
		Where("cocktail_id = ?", id).
		Find(&ingredients).Error

	if err != nil {
		return []domain.CocktailIngredient{}, err
	}

	return ingredients, nil
}
