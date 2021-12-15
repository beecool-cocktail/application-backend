package usecase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type cocktailUsecase struct {
	cocktailMySQLRepo domain.CocktailMySQLRepository
}

// NewDietUsecase ...
func NewCocktailUsecase(cocktailMySQLRepo domain.CocktailMySQLRepository) domain.CocktailUsecase {
	return &cocktailUsecase{
		cocktailMySQLRepo: cocktailMySQLRepo,
	}
}

func (c *cocktailUsecase) GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination domain.PaginationUsecase) (*[]domain.Cocktail, int64, error) {
	sortByDir := make(map[string]sortbydir.SortByDir)
	for sort, dir := range pagination.SortByDir {
		sortByDir[sort] = sortbydir.ParseSortByDirByInt(dir)
	}

	cocktails, total, err := c.cocktailMySQLRepo.GetAllWithFilter(ctx, filter, domain.PaginationMySQLRepository{
		Page: pagination.Page,
		PageSize: pagination.PageSize,
		SortByDir:sortByDir,
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return nil, 0, domain.ErrCocktailNotFound
	} else if err != nil {
		logrus.Error(err)
		return nil, 0, err
	}

	return cocktails, total, nil
}