package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
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

func (s *cocktailStepMySQLRepository) QueryByCocktailId(ctx context.Context, id int64) ([]domain.CocktailStep, error) {

	var steps []domain.CocktailStep
	order := sortbydir.MakeSortAndDir("step_number", sortbydir.ParseStringBySortByDir(sortbydir.ASC))
	err := s.db.Select("id", "step_description").
		Where("cocktail_id = ?", id).
		Order(order).
		Find(&steps).Error

	if err != nil {
		return []domain.CocktailStep{}, err
	}

	return steps, nil
}
