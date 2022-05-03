package usecase

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"gorm.io/gorm"
)

type favoriteCocktailUsecase struct {
	favoriteCocktailMySQL  domain.FavoriteCocktailMySQLRepository
	cocktailMySQL          domain.CocktailMySQLRepository
	cocktailPhotoMySQLRepo domain.CocktailPhotoMySQLRepository
	userMySQLRepo          domain.UserMySQLRepository
	userRedisRepo          domain.UserRedisRepository
	transactionRepo        domain.DBTransactionRepository
}

func NewFavoriteCocktailUsecase(
	favoriteCocktailMySQL domain.FavoriteCocktailMySQLRepository,
	cocktailMySQL domain.CocktailMySQLRepository,
	cocktailPhotoMySQLRepo domain.CocktailPhotoMySQLRepository,
	userMySQLRepo domain.UserMySQLRepository,
	userRedisRepo domain.UserRedisRepository,
	transactionRepo domain.DBTransactionRepository) domain.FavoriteCocktailUsecase {
	return &favoriteCocktailUsecase{
		favoriteCocktailMySQL:  favoriteCocktailMySQL,
		cocktailMySQL:          cocktailMySQL,
		cocktailPhotoMySQLRepo: cocktailPhotoMySQLRepo,
		userMySQLRepo:          userMySQLRepo,
		userRedisRepo:          userRedisRepo,
		transactionRepo:        transactionRepo,
	}
}

func (f *favoriteCocktailUsecase) fillFavoriteCocktailList(ctx context.Context,
	cocktails []domain.APIFavoriteCocktail) ([]domain.APIFavoriteCocktail, error) {

	var apiFavoriteCocktails []domain.APIFavoriteCocktail

	for _, favoriteCocktail := range cocktails {
		cocktail, err := f.cocktailMySQL.QueryByCocktailID(ctx, favoriteCocktail.CocktailID)
		if err != nil && err == gorm.ErrRecordNotFound {
			// origin article was deleted
			deleteErr := f.favoriteCocktailMySQL.Delete(ctx, favoriteCocktail.CocktailID, favoriteCocktail.UserID)
			if deleteErr != nil {
				return []domain.APIFavoriteCocktail{}, err
			}
			continue
		} else if err != nil {
			return []domain.APIFavoriteCocktail{}, err
		}

		photo, err := f.cocktailPhotoMySQLRepo.QueryCoverPhotoByCocktailId(ctx, cocktail.CocktailID)
		if err != nil {
			return []domain.APIFavoriteCocktail{}, err
		}

		userName, err := f.userRedisRepo.QueryUserNameByID(ctx, cocktail.UserID)
		if err != nil {
			return []domain.APIFavoriteCocktail{}, err
		}

		out := domain.APIFavoriteCocktail{
			CocktailID: cocktail.CocktailID,
			UserID:     cocktail.UserID,
			UserName:   userName,
			Title:      cocktail.Title,
			CoverPhoto: photo,
		}

		apiFavoriteCocktails = append(apiFavoriteCocktails, out)
	}

	return apiFavoriteCocktails, nil
}

func (f *favoriteCocktailUsecase) Store(ctx context.Context, c *domain.FavoriteCocktail) error {

	user, err := f.userMySQLRepo.QueryById(ctx, c.UserID)
	if err != nil {
		return err
	}

	if err := f.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		err := f.favoriteCocktailMySQL.StoreTx(ctx, tx, c)
		if err != nil {
			return err
		}

		user.NumberOfCollection = user.NumberOfCollection + 1
		_, err = f.userMySQLRepo.UpdateNumberOfNumberOfCollectionTx(ctx, tx, user)
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

func (f *favoriteCocktailUsecase) Delete(ctx context.Context, cocktailID, userID int64) error {
	user, err := f.userMySQLRepo.QueryById(ctx, userID)
	if err != nil {
		return err
	}

	if err := f.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		err := f.favoriteCocktailMySQL.DeleteTx(ctx, tx, cocktailID, userID)
		if err != nil {
			return err
		}

		user.NumberOfCollection = user.NumberOfCollection - 1
		_, err = f.userMySQLRepo.UpdateNumberOfNumberOfCollectionTx(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
