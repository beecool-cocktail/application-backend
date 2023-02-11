package usecase

import (
	"context"
	"fmt"
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

		cocktail, err := f.cocktailMySQL.QueryByCocktailID(ctx, favoriteCocktail.CocktailID)
		if err != nil {
			return []domain.APIFavoriteCocktail{}, err
		}

		photo, err := f.cocktailPhotoMySQLRepo.QueryCoverPhotoByCocktailId(ctx, favoriteCocktail.CocktailID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return []domain.APIFavoriteCocktail{}, err
		}

		userName, err := f.userRedisRepo.QueryUserNameByID(ctx, favoriteCocktail.UserID)
		if err != nil {
			return []domain.APIFavoriteCocktail{}, err
		}

		out := domain.APIFavoriteCocktail{
			CocktailID:    favoriteCocktail.CocktailID,
			UserID:        favoriteCocktail.UserID,
			UserName:      userName,
			Title:         cocktail.Title,
			CoverPhoto:    photo,
			CollectedDate: favoriteCocktail.CollectedDate,
		}

		apiFavoriteCocktails = append(apiFavoriteCocktails, out)
	}

	return apiFavoriteCocktails, nil
}

func (f *favoriteCocktailUsecase) fillCollectionStatusInList(ctx context.Context, cocktails []domain.APIFavoriteCocktail,
	userID int64) ([]domain.APIFavoriteCocktail, error) {

	var apiCocktails []domain.APIFavoriteCocktail
	favoriteCocktails, _, err := f.favoriteCocktailMySQL.QueryByUserID(ctx, userID, domain.PaginationMySQLRepository{})
	if err != nil {
		return []domain.APIFavoriteCocktail{}, err
	}

	favoriteCocktailsMap := make(map[int64]bool)
	for _, favoriteCocktail := range favoriteCocktails {
		favoriteCocktailsMap[favoriteCocktail.CocktailID] = true
	}

	for _, cocktail := range cocktails {
		if _, ok := favoriteCocktailsMap[cocktail.CocktailID]; ok {
			cocktail.IsCollected = true
		} else {
			cocktail.IsCollected = false
		}

		apiCocktails = append(apiCocktails, cocktail)
	}

	return apiCocktails, nil
}

func (f *favoriteCocktailUsecase) Store(ctx context.Context, c *domain.FavoriteCocktail) error {

	if err := f.transactionRepo.Transaction(func(i interface{}) (err error) {
		tx := i.(*gorm.DB)

		cocktailID := strconv.FormatInt(c.CocktailID, 10)
		retryInterval := 100 * time.Millisecond
		retryTimes := 5
		lock, err := f.cocktailRedisRepo.GetCocktailCollectionNumberLock(ctx, cocktailID, time.Second, retryInterval,
			retryTimes)
		if err != nil {
			return err
		}

		defer func() {
			if deferError := f.cocktailRedisRepo.ReleaseCocktailCollectionNumberLock(ctx, lock); deferError != nil {
				err = fmt.Errorf("error is: %w, and defer error is %s", err, deferError)
			}
		}()

		err = f.favoriteCocktailMySQL.StoreTx(ctx, tx, c)
		if err != nil {
			return err
		}

		_, err = f.cocktailMySQL.IncreaseNumberOfCollectionTx(ctx, tx, c.CocktailID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (f *favoriteCocktailUsecase) QueryByUserID(ctx context.Context, targetUserID int64,
	pagination domain.PaginationUsecase, queryUserID int64) ([]domain.APIFavoriteCocktail, int64, error) {

	var apiFavoriteCocktails []domain.APIFavoriteCocktail

	sortByDir := make(map[string]sortbydir.SortByDir)
	for sort, dir := range pagination.SortByDir {
		sortByDir[sort] = sortbydir.ParseSortByDirByInt(dir)
	}

	sortByDir["created_date"] = sortbydir.ParseSortByDirByInt(1)

	cocktails, total, err := f.favoriteCocktailMySQL.QueryByUserID(ctx, targetUserID,
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
			CocktailID:    cocktail.CocktailID,
			UserID:        cocktail.UserID,
			CollectedDate: util.GetFormatTime(cocktail.CreatedDate, "UTC"),
		}
		apiFavoriteCocktails = append(apiFavoriteCocktails, out)
	}

	apiFavoriteCocktail, err := f.fillFavoriteCocktailList(ctx, apiFavoriteCocktails)
	if err != nil {
		return []domain.APIFavoriteCocktail{}, total, err
	}

	if queryUserID != 0 {
		apiFavoriteCocktail, err = f.fillCollectionStatusInList(ctx, apiFavoriteCocktail, queryUserID)
		if err != nil {
			return []domain.APIFavoriteCocktail{}, 0, err
		}
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

	if err := f.transactionRepo.Transaction(func(i interface{}) (err error) {
		tx := i.(*gorm.DB)

		cocktailIDKey := strconv.FormatInt(cocktailID, 10)
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

		err = f.favoriteCocktailMySQL.DeleteTx(ctx, tx, cocktailID, userID)
		if err != nil {
			return err
		}

		_, err = f.cocktailMySQL.DecreaseNumberOfCollectionTx(ctx, tx, cocktailID)
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
