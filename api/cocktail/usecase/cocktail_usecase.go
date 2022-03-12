package usecase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cocktailUsecase struct {
	cocktailMySQLRepo           domain.CocktailMySQLRepository
	cocktailFileRepo            domain.CocktailFileRepository
	cocktailPhotoMySQLRepo      domain.CocktailPhotoMySQLRepository
	cocktailIngredientMySQLRepo domain.CocktailIngredientMySQLRepository
	cocktailStepMySQLRepo       domain.CocktailStepMySQLRepository
	transactionRepo             domain.DBTransactionRepository
}

// NewDietUsecase ...
func NewCocktailUsecase(
	cocktailMySQLRepo domain.CocktailMySQLRepository,
	cocktailFileRepo domain.CocktailFileRepository,
	cocktailPhotoMySQLRepo domain.CocktailPhotoMySQLRepository,
	cocktailIngredientMySQLRepo domain.CocktailIngredientMySQLRepository,
	cocktailStepMySQLRepo domain.CocktailStepMySQLRepository,
	transactionRepo domain.DBTransactionRepository) domain.CocktailUsecase {
	return &cocktailUsecase{
		cocktailMySQLRepo:           cocktailMySQLRepo,
		cocktailFileRepo:            cocktailFileRepo,
		cocktailPhotoMySQLRepo:      cocktailPhotoMySQLRepo,
		cocktailIngredientMySQLRepo: cocktailIngredientMySQLRepo,
		cocktailStepMySQLRepo:       cocktailStepMySQLRepo,
		transactionRepo:             transactionRepo,
	}
}

func (c *cocktailUsecase) fillCocktailCoverPhoto(ctx context.Context, cocktails []domain.APICocktail) ([]domain.APICocktail, error) {

	var apiCocktails []domain.APICocktail

	for _, cocktail := range cocktails {
		path, err := c.cocktailPhotoMySQLRepo.QueryCoverPhotoByCocktailId(ctx, cocktail.CocktailID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiCocktails = append(apiCocktails, cocktail)
			continue
		} else if err != nil {
			return []domain.APICocktail{}, err
		}
		cocktail.CoverPhoto = path
		apiCocktails = append(apiCocktails, cocktail)
	}

	return apiCocktails, nil
}

func (c *cocktailUsecase) fillCocktailDetails(ctx context.Context, cocktail domain.APICocktail) (domain.APICocktail, error) {

	paths, err := c.cocktailPhotoMySQLRepo.QueryPhotosByCocktailId(ctx, cocktail.CocktailID)
	if err != nil {
		return domain.APICocktail{}, err
	}
	for _, path := range paths {
		cocktail.Photos = append(cocktail.Photos, path)
	}

	ingredients, err := c.cocktailIngredientMySQLRepo.QueryByCocktailId(ctx, cocktail.CocktailID)
	if err != nil {
		return domain.APICocktail{}, err
	}
	cocktail.Ingredients = ingredients

	steps, err := c.cocktailStepMySQLRepo.QueryByCocktailId(ctx, cocktail.CocktailID)
	if err != nil {
		return domain.APICocktail{}, err
	}
	cocktail.Steps = steps

	return cocktail, nil
}

func (c *cocktailUsecase) GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination domain.PaginationUsecase) ([]domain.APICocktail, int64, error) {
	sortByDir := make(map[string]sortbydir.SortByDir)
	for sort, dir := range pagination.SortByDir {
		sortByDir[sort] = sortbydir.ParseSortByDirByInt(dir)
	}

	sortByDir["created_date"] = sortbydir.ParseSortByDirByInt(1)

	cocktails, total, err := c.cocktailMySQLRepo.GetAllWithFilter(ctx, filter, domain.PaginationMySQLRepository{
		Page:      pagination.Page,
		PageSize:  pagination.PageSize,
		SortByDir: sortByDir,
	})
	if err != nil {
		return []domain.APICocktail{}, 0, err
	}

	var apiCocktails []domain.APICocktail
	for _, cocktail := range cocktails {
		out := domain.APICocktail{
			CocktailID:  cocktail.CocktailID,
			UserID:      cocktail.UserID,
			Title:       cocktail.Title,
			Description: cocktail.Description,
			CreatedDate: util.GetFormatTime(cocktail.CreatedDate, "UTC"),
		}
		apiCocktails = append(apiCocktails, out)
	}

	apiCocktails, err = c.fillCocktailCoverPhoto(ctx, apiCocktails)
	if err != nil {
		return []domain.APICocktail{}, 0, err
	}

	return apiCocktails, total, nil
}

func (c *cocktailUsecase) QueryByCocktailID(ctx context.Context, id int64) (domain.APICocktail, error) {

	cocktail, err := c.cocktailMySQLRepo.QueryByCocktailID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.APICocktail{}, domain.ErrCocktailNotFound
	} else if err != nil {
		return domain.APICocktail{}, err
	}

	apiCocktail := domain.APICocktail{
		CocktailID:  cocktail.CocktailID,
		UserID:      cocktail.UserID,
		Title:       cocktail.Title,
		Description: cocktail.Description,
		CreatedDate: util.GetFormatTime(cocktail.CreatedDate, "UTC"),
	}

	apiCocktail, err = c.fillCocktailDetails(ctx, apiCocktail)
	if err != nil {
		return domain.APICocktail{}, err
	}

	return apiCocktail, nil
}

func (c *cocktailUsecase) Store(ctx context.Context, co *domain.Cocktail, ingredients []domain.CocktailIngredient,
	steps []domain.CocktailStep, images []domain.CocktailImage) error {

	//Todo move to config
	savePath := "static/images/"
	urlPath := "static/"

	newCocktailID := util.GetID(util.IdGenerator)

	err := c.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		co.CocktailID = newCocktailID
		err := c.cocktailMySQLRepo.StoreTx(ctx, tx, co)
		if err != nil {
			return err
		}

		for _, image := range images {

			newFileName := uuid.New().String()
			image.Name = newFileName

			if !util.ValidateImageType(image.Type) {
				return domain.ErrCodeFileTypeIllegal
			}

			image.Destination = savePath + newFileName
			err := c.cocktailFileRepo.SaveAsWebp(ctx, &image)
			if err != nil {
				return err
			}

			image.Destination = urlPath + newFileName + ".webp"
			err = c.cocktailPhotoMySQLRepo.StoreTx(ctx, tx,
				&domain.CocktailPhoto{
					CocktailID:   newCocktailID,
					Photo:        image.Destination,
					IsCoverPhoto: image.IsCoverPhoto,
				})
			if err != nil {
				return err
			}
		}

		for _, ingredient := range ingredients {
			ingredient.CocktailID = newCocktailID
			err = c.cocktailIngredientMySQLRepo.StoreTx(ctx, tx, &ingredient)
			if err != nil {
				return err
			}
		}

		for _, step := range steps {
			step.CocktailID = newCocktailID
			err = c.cocktailStepMySQLRepo.StoreTx(ctx, tx, &step)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
