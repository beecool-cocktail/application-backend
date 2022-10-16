package command

import (
	"context"
	"fmt"
	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type favoriteCocktail struct {
	favoriteCocktailMySQL domain.FavoriteCocktailMySQLRepository
	cocktailMySQLRepo     domain.CocktailMySQLRepository
	cocktailRedisRepo     domain.CocktailRedisRepository
	transactionRepo       domain.DBTransactionRepository
}

func NewFavoriteCocktailOperator(
	favoriteCocktailMySQL domain.FavoriteCocktailMySQLRepository,
	cocktailMySQLRepo domain.CocktailMySQLRepository,
	cocktailRedisRepo domain.CocktailRedisRepository,
	transactionRepo domain.DBTransactionRepository) Operator {
	return &favoriteCocktail{
		favoriteCocktailMySQL: favoriteCocktailMySQL,
		cocktailMySQLRepo:     cocktailMySQLRepo,
		cocktailRedisRepo:     cocktailRedisRepo,
		transactionRepo:       transactionRepo,
	}
}

func (f *favoriteCocktail) Undo(ctx context.Context, command *domain.Command) error {

	if err := f.transactionRepo.Transaction(func(i interface{}) (err error) {
		tx := i.(*gorm.DB)

		cocktailID := command.Type.Delete.TargetID.(float64)
		userID := command.Type.Delete.OperatorID.(float64)

		cocktailIDKey := strconv.FormatFloat(cocktailID, 'f', 0, 64)
		retryInterval := 100 * time.Millisecond
		retryTimes := 5
		lock, err := f.cocktailRedisRepo.GetCocktailCollectionNumberLock(ctx, cocktailIDKey, time.Second, retryInterval,
			retryTimes)
		if err != nil {
			return err
		}

		defer func() {
			if deferError := f.cocktailRedisRepo.ReleaseCocktailCollectionNumberLock(ctx, lock); deferError != nil {
				err = fmt.Errorf("error is: %w, and defer error is %s", err, deferError)
			}
		}()

		err = f.favoriteCocktailMySQL.StoreTx(ctx, tx, &domain.FavoriteCocktail{
			CocktailID: int64(cocktailID),
			UserID:     int64(userID),
		})
		if err != nil {
			return err
		}

		_, err = f.cocktailMySQLRepo.IncreaseNumberOfCollectionTx(ctx, tx, int64(cocktailID))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
