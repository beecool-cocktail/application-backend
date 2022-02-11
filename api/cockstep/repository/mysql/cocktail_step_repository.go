package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
)

type cocktailStepMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLCocktailStepRepository(db *gorm.DB) domain.CocktailStepMySQLRepository {
	return &cocktailStepMySQLRepository{db}
}

func (s *cocktailStepMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailStep) error {

	err := tx.Select("cocktail_id", "step_number", "step_description").Create(c).Error

	if err != nil {
		return err
	}

	return nil
}

