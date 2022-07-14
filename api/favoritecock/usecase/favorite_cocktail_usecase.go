package usecase

import (
	"context"
	"github.com/beecool-cocktail/application-backend/command"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"github.com/beecool-cocktail/application-backend/util"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type favoriteCocktailUsecase struct {
	favoriteCocktailMySQL  domain.FavoriteCocktailMySQLRepository
	cocktailMySQL          domain.CocktailMySQLRepository
	cocktailRedisRepo      domain.CocktailRedisRepository
	cocktailPhotoMySQLRepo domain.CocktailPhotoMySQLRepository
	userMySQLRepo          domain.UserMySQLRepository
	userRedisRepo          domain.UserRedisRepository
	commandRedisRepo       domain.CommandRedisRepository
	transactionRepo        domain.DBTransactionRepository
}

func NewFavoriteCocktailUsecase(
	favoriteCocktailMySQL domain.FavoriteCocktailMySQLRepository,
	cocktailMySQL domain.CocktailMySQLRepository,
	cocktailRedisRepo domain.CocktailRedisRepository,
	cocktailPhotoMySQLRepo domain.CocktailPhotoMySQLRepository,
	userMySQLRepo domain.UserMySQLRepository,
	userRedisRepo domain.UserRedisRepository,
	commandRedisRepo domain.CommandRedisRepository,
	transactionRepo domain.DBTransactionRepository) domain.FavoriteCocktailUsecase {
	return &favoriteCocktailUsecase{
		favoriteCocktailMySQL:  favoriteCocktailMySQL,
		cocktailMySQL:          cocktailMySQL,
		cocktailRedisRepo:      cocktailRedisRepo,
		cocktailPhotoMySQLRepo: cocktailPhotoMySQLRepo,
		userMySQLRepo:          userMySQLRepo,
		userRedisRepo:          userRedisRepo,
		commandRedisRepo:       commandRedisRepo,
		transactionRepo:        transactionRepo,
	}
}

func (f *favoriteCocktailUsecase) fillFavoriteCocktailList(ctx context.Context,
	cocktails []domain.APIFavoriteCocktail) ([]domain.APIFavoriteCocktail, error) {

	var apiFavoriteCocktails []domain.APIFavoriteCocktail

	for _, favoriteCocktail := range cocktails {

		photo, err := f.cocktailPhotoMySQLRepo.QueryCoverPhotoByCocktailId(ctx, favoriteCocktail.CocktailID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return []domain.APIFavoriteCocktail{}, err
		}

		userName, err := f.userRedisRepo.QueryUserNameByID(ctx, favoriteCocktail.UserID)
		if err != nil {
			return []domain.APIFavoriteCocktail{}, err
		}

		out := domain.APIFavoriteCocktail{
			CocktailID: favoriteCocktail.CocktailID,
			UserID:     favoriteCocktail.UserID,
			UserName:   userName,
			Title:      favoriteCocktail.Title,
			CoverPhoto: photo,
		}

		apiFavoriteCocktails = append(apiFavoriteCocktails, out)
	}

	return apiFavoriteCocktails, nil
}

func (f *favoriteCocktailUsecase) Store(ctx context.Context, c *domain.FavoriteCocktail) error {

	if err := f.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		err := f.favoriteCocktailMySQL.StoreTx(ctx, tx, c)
		if err != nil {
			return err
		}

		err = f.cocktailRedisRepo.IncreaseCollectionNumbers(ctx, &domain.CocktailCollection{
			CocktailID:       c.CocktailID,
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

func (f *favoriteCocktailUsecase) QueryByUserID(ctx context.Context, id int64,
	pagination domain.PaginationUsecase) ([]domain.APIFavoriteCocktail, int64, error) {

	var apiFavoriteCocktails []domain.APIFavoriteCocktail

	sortByDir := make(map[string]sortbydir.SortByDir)
	for sort, dir := range pagination.SortByDir {
		sortByDir[sort] = sortbydir.ParseSortByDirByInt(dir)
	}

	sortByDir["created_date"] = sortbydir.ParseSortByDirByInt(1)

	cocktails, total, err := f.favoriteCocktailMySQL.QueryByUserID(ctx, id,
		domain.PaginationMySQLRepository{
			Page:      pagination.Page,
			PageSize:  pagination.PageSize,
			SortByDir: sortByDir,
		})

	if err != nil {
		return []domain.APIFavoriteCocktail{}, total, err
	}

	for _, cocktail := range cocktails {
		out := domain.APIFavoriteCocktail{
			CocktailID: cocktail.CocktailID,
			UserID:     cocktail.UserID,
		}
		apiFavoriteCocktails = append(apiFavoriteCocktails, out)
	}

	apiFavoriteCocktail, err := f.fillFavoriteCocktailList(ctx, apiFavoriteCocktails)
	if err != nil {
		return []domain.APIFavoriteCocktail{}, total, err
	}

	return apiFavoriteCocktail, total, nil
}

func (f *favoriteCocktailUsecase) QueryCountsByUserID(ctx context.Context, id int64) (int64, error) {

	total, err := f.favoriteCocktailMySQL.QueryCountsByUserID(ctx, id)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (f *favoriteCocktailUsecase) Delete(ctx context.Context, cocktailID, userID int64) (string, error) {

	commandID := strconv.FormatInt(util.GetID(util.IdGenerator), 10)

	if err := f.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		err := f.favoriteCocktailMySQL.DeleteTx(ctx, tx, cocktailID, userID)
		if err != nil {
			return err
		}

		err = f.cocktailRedisRepo.DecreaseCollectionNumbers(ctx, &domain.CocktailCollection{
			CocktailID:       cocktailID,
			CollectionCounts: 1,
		})
		if err != nil {
			return err
		}

		err = f.commandRedisRepo.Store(ctx, &domain.Command{
			ID:   commandID,
			Name: command.FavoriteCocktailDelete,
			Type: domain.CommandType{
				Delete: domain.DeleteCommand{
					OperatorID: userID,
					TargetID:   cocktailID,
				},
			},
			ExpireTime: time.Minute * time.Duration(3),
		})

		return nil
	}); err != nil {
		return commandID, err
	}

	return commandID, nil
}
