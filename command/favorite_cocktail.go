package command

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
)

type favoriteCocktail struct {
	favoriteCocktailMySQL domain.FavoriteCocktailMySQLRepository
	cocktailRedisRepo     domain.CocktailRedisRepository
	transactionRepo       domain.DBTransactionRepository
}

func NewFavoriteCocktailOperator(
	favoriteCocktailMySQL domain.FavoriteCocktailMySQLRepository,
	cocktailRedisRepo domain.CocktailRedisRepository,
	transactionRepo domain.DBTransactionRepository) Operator {
	return &favoriteCocktail{
		favoriteCocktailMySQL: favoriteCocktailMySQL,
		cocktailRedisRepo:     cocktailRedisRepo,
		transactionRepo:       transactionRepo,
	}
}

func (f *favoriteCocktail) Undo(ctx context.Context, command *domain.Command) error {

	if err := f.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		cocktailID := command.Type.Delete.TargetID.(float64)
		userID := command.Type.Delete.OperatorID.(float64)

		err := f.favoriteCocktailMySQL.StoreTx(ctx, tx, &domain.FavoriteCocktail{
			CocktailID: int64(cocktailID),
			UserID:     int64(userID),
		})
		if err != nil {
			return err
		}

		err = f.cocktailRedisRepo.IncreaseCollectionNumbers(ctx, &domain.CocktailCollection{
			CocktailID:       int64(cocktailID),
			CollectionCounts: 1,
		})
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
