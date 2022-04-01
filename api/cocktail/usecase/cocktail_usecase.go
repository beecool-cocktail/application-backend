package usecase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/cockarticletype"
	"github.com/beecool-cocktail/application-backend/enum/httpaction"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type cocktailUsecase struct {
	service                     *service.Service
	cocktailMySQLRepo           domain.CocktailMySQLRepository
	cocktailFileRepo            domain.CocktailFileRepository
	cocktailPhotoMySQLRepo      domain.CocktailPhotoMySQLRepository
	cocktailIngredientMySQLRepo domain.CocktailIngredientMySQLRepository
	cocktailStepMySQLRepo       domain.CocktailStepMySQLRepository
	userMySQLRepo               domain.UserMySQLRepository
	transactionRepo             domain.DBTransactionRepository
}

// NewDietUsecase ...
func NewCocktailUsecase(
	s *service.Service,
	cocktailMySQLRepo domain.CocktailMySQLRepository,
	cocktailFileRepo domain.CocktailFileRepository,
	cocktailPhotoMySQLRepo domain.CocktailPhotoMySQLRepository,
	cocktailIngredientMySQLRepo domain.CocktailIngredientMySQLRepository,
	cocktailStepMySQLRepo domain.CocktailStepMySQLRepository,
	userMySQLRepo domain.UserMySQLRepository,
	transactionRepo domain.DBTransactionRepository) domain.CocktailUsecase {
	return &cocktailUsecase{
		service:                     s,
		cocktailMySQLRepo:           cocktailMySQLRepo,
		cocktailFileRepo:            cocktailFileRepo,
		cocktailPhotoMySQLRepo:      cocktailPhotoMySQLRepo,
		cocktailIngredientMySQLRepo: cocktailIngredientMySQLRepo,
		cocktailStepMySQLRepo:       cocktailStepMySQLRepo,
		userMySQLRepo:               userMySQLRepo,
		transactionRepo:             transactionRepo,
	}
}

func (c *cocktailUsecase) fillCocktailList(ctx context.Context, cocktails []domain.APICocktail) ([]domain.APICocktail, error) {

	var apiCocktails []domain.APICocktail

	for _, cocktail := range cocktails {
		photos, err := c.cocktailPhotoMySQLRepo.QueryPhotosByCocktailId(ctx, cocktail.CocktailID)
		if err != nil {
			return []domain.APICocktail{}, err
		}
		cocktail.Photos = photos
		for _, photo := range photos {
			if photo.IsCoverPhoto == true {
				cocktail.CoverPhoto = photo
			}
		}

		ingredients, err := c.cocktailIngredientMySQLRepo.QueryByCocktailId(ctx, cocktail.CocktailID)
		if err != nil {
			return []domain.APICocktail{}, err
		}
		cocktail.Ingredients = ingredients

		user, err := c.userMySQLRepo.QueryById(ctx, cocktail.UserID)
		if err != nil {
			return []domain.APICocktail{}, err
		}
		cocktail.UserName = user.Name

		apiCocktails = append(apiCocktails, cocktail)
	}

	return apiCocktails, nil
}

func (c *cocktailUsecase) fillCocktailDetails(ctx context.Context, cocktail domain.APICocktail) (domain.APICocktail, error) {

	photos, err := c.cocktailPhotoMySQLRepo.QueryPhotosByCocktailId(ctx, cocktail.CocktailID)
	if err != nil {
		return domain.APICocktail{}, err
	}
	cocktail.Photos = photos

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

	user, err := c.userMySQLRepo.QueryById(ctx, cocktail.UserID)
	if err != nil {
		return domain.APICocktail{}, err
	}
	cocktail.UserName = user.Name

	return cocktail, nil
}

func (c *cocktailUsecase) getNeedDeletedPhoto(oldPhotos []domain.CocktailPhoto,
	newImages []domain.CocktailImage) []int64 {

	var deletedPhotoID []int64
	for _, oldPhoto := range oldPhotos {
		var needDeleted = true
		for _, newImage := range newImages {
			if oldPhoto.ID == newImage.ImageID {
				needDeleted = false
			}
		}
		if needDeleted {
			deletedPhotoID = append(deletedPhotoID, oldPhoto.ID)
		}
	}

	return deletedPhotoID
}

func (c *cocktailUsecase) getAction(id int64, data string) (httpaction.HttpAction, error) {

	if id > 0 && data != "" {
		return httpaction.Edit, nil
	} else if id > 0 && data == "" {
		return httpaction.Keep, nil
	} else if id == 0 && data != "" {
		return httpaction.Add, nil
	} else {
		return httpaction.Keep, domain.ErrCanNotSpecifyHttpAction
	}
}

func (c *cocktailUsecase) addPhoto(ctx context.Context, tx *gorm.DB, image *domain.CocktailImage) error {
	savePath := c.service.Configure.Others.File.Image.PathInDB
	urlPath := c.service.Configure.Others.File.Image.PathInURL

	newFileName := uuid.New().String()
	image.Name = newFileName

	if !util.ValidateImageType(image.Type) {
		return domain.ErrCodeFileTypeIllegal
	}

	image.Destination = savePath + newFileName
	err := c.cocktailFileRepo.SaveAsWebp(ctx, image)
	if err != nil {
		return err
	}

	image.Destination = urlPath + newFileName + ".webp"
	err = c.cocktailPhotoMySQLRepo.StoreTx(ctx, tx,
		&domain.CocktailPhoto{
			CocktailID:   image.CocktailID,
			Photo:        image.Destination,
			IsCoverPhoto: image.IsCoverPhoto,
		})
	if err != nil {
		return err
	}

	return nil
}

func (c *cocktailUsecase) editPhoto(ctx context.Context, tx *gorm.DB, image *domain.CocktailImage) error {
	savePath := c.service.Configure.Others.File.Image.PathInDB
	urlPath := c.service.Configure.Others.File.Image.PathInURL

	photo, err := c.cocktailPhotoMySQLRepo.QueryPhotoById(ctx, image.ImageID)
	if err != nil {
		return err
	}

	fileName, err := util.GetFileNameByPath(photo.Photo)
	if err != nil {
		return err
	}

	image.Destination = savePath + fileName
	err = c.cocktailFileRepo.UpdateAsWebp(ctx, image)
	if err != nil {
		return err
	}

	//this file already have type
	image.Destination = urlPath + fileName
	_, err = c.cocktailPhotoMySQLRepo.UpdateTx(ctx, tx,
		&domain.CocktailPhoto{
			ID:           image.ImageID,
			Photo:        image.Destination,
			IsCoverPhoto: image.IsCoverPhoto,
		})
	if err != nil {
		return err
	}

	return nil
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

	apiCocktails, err = c.fillCocktailList(ctx, apiCocktails)
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

func (c *cocktailUsecase) QueryDraftByCocktailID(ctx context.Context, cocktailID, userID int64) (domain.APICocktail, error) {

	cocktail, err := c.cocktailMySQLRepo.QueryByCocktailID(ctx, cocktailID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.APICocktail{}, domain.ErrCocktailNotFound
	} else if err != nil {
		return domain.APICocktail{}, err
	}

	if cocktail.UserID != userID {
		return domain.APICocktail{}, domain.ErrItemDoesNotBelongToUser
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
	steps []domain.CocktailStep, images []domain.CocktailImage, userID int64) error {

	savePath := c.service.Configure.Others.File.Image.PathInDB
	urlPath := c.service.Configure.Others.File.Image.PathInURL

	newCocktailID := util.GetID(util.IdGenerator)

	user, err := c.userMySQLRepo.QueryById(ctx, userID)
	if err != nil {
		return err
	}

	err = c.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		if co.Category == cockarticletype.Draft.Int() {
			if user.NumberOfDraft == 30 {
				return domain.ErrorCocktailDraftIsMaximum
			}
			numberOfDraft := user.NumberOfDraft + 1
			_, err = c.userMySQLRepo.UpdateNumberOfDraftTx(ctx, tx, &domain.User{
				ID:            user.ID,
				NumberOfDraft: numberOfDraft,
			})
			if err != nil {
				return err
			}
		} else {
			numberOfPost := user.NumberOfPost + 1
			_, err = c.userMySQLRepo.UpdateNumberOfPostTx(ctx, tx, &domain.User{
				ID:           user.ID,
				NumberOfPost: numberOfPost,
			})
			if err != nil {
				return err
			}
		}

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

func (c *cocktailUsecase) Update(ctx context.Context, co *domain.Cocktail, ingredients []domain.CocktailIngredient,
	steps []domain.CocktailStep, images []domain.CocktailImage, userID int64) error {

	if co.CocktailID <= 0 {
		return domain.ErrParameterIllegal
	}

	cocktail, err := c.cocktailMySQLRepo.QueryByCocktailID(ctx, co.CocktailID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrCocktailNotFound
	} else if err != nil {
		return err
	}

	// cocktail isn't belong to user
	if cocktail.UserID != userID {
		return domain.ErrItemDoesNotBelongToUser
	}

	// can't update non-draft article
	if cocktail.Category != cockarticletype.Draft.Int() {
		return domain.ErrPermissionDenied
	}

	oldPhoto, err := c.cocktailPhotoMySQLRepo.QueryPhotosByCocktailId(ctx, co.CocktailID)
	if err != nil {
		return err
	}

	deletedPhoto := c.getNeedDeletedPhoto(oldPhoto, images)

	if err := c.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		// update cocktail basic info
		_, err := c.cocktailMySQLRepo.UpdateTx(ctx, tx, co)
		if err != nil {
			return err
		}

		// update cocktail photo
		for _, image := range images {
			action, err := c.getAction(image.ImageID, image.Data)
			if err != nil {
				return err
			}

			logrus.Debugf("update photo id: %d", image.ImageID)
			logrus.Debugf("update photo action: %s", action.String())

			if action == httpaction.Add {
				err = c.addPhoto(ctx, tx, &image)
				if err != nil {
					return err
				}
			} else if action == httpaction.Edit {
				err = c.editPhoto(ctx, tx, &image)
				if err != nil {
					return err
				}
			} else {
				//keep photo
			}
		}

		for _, id := range deletedPhoto {
			err := c.cocktailPhotoMySQLRepo.DeleteByIDTx(ctx, tx, id)
			if err != nil {
				return err
			}
		}

		//update cocktail ingredient
		err = c.cocktailIngredientMySQLRepo.DeleteByCocktailIDTx(ctx, tx, co.CocktailID)
		if err != nil {
			return err
		}

		for _, ingredient := range ingredients {
			err = c.cocktailIngredientMySQLRepo.StoreTx(ctx, tx, &ingredient)
			if err != nil {
				return err
			}
		}

		//update cocktail step
		err = c.cocktailStepMySQLRepo.DeleteByCocktailIDTx(ctx, tx, co.CocktailID)
		if err != nil {
			return err
		}

		for _, step := range steps {
			err = c.cocktailStepMySQLRepo.StoreTx(ctx, tx, &step)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *cocktailUsecase) Delete(ctx context.Context, cocktailID, userID int64) error {

	cocktail, err := c.cocktailMySQLRepo.QueryByCocktailID(ctx, cocktailID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrCocktailNotFound
	} else if err != nil {
		return err
	}

	if cocktail.UserID != userID {
		return domain.ErrItemDoesNotBelongToUser
	}

	user, err := c.userMySQLRepo.QueryById(ctx, userID)
	if err != nil {
		return err
	}

	if err := c.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		if cocktail.Category == cockarticletype.Draft.Int() {
			if user.NumberOfDraft == 30 {
				return domain.ErrorCocktailDraftIsMaximum
			}
			numberOfDraft := user.NumberOfDraft - 1
			_, err = c.userMySQLRepo.UpdateNumberOfDraftTx(ctx, tx, &domain.User{
				ID:            user.ID,
				NumberOfDraft: numberOfDraft,
			})
			if err != nil {
				return err
			}
		} else {
			numberOfPost := user.NumberOfPost - 1
			_, err = c.userMySQLRepo.UpdateNumberOfPostTx(ctx, tx, &domain.User{
				ID:           user.ID,
				NumberOfPost: numberOfPost,
			})
			if err != nil {
				return err
			}
		}

		err := c.cocktailMySQLRepo.DeleteTx(ctx, tx, cocktailID)
		if err != nil {
			return err
		}

		err = c.cocktailIngredientMySQLRepo.DeleteByCocktailIDTx(ctx, tx, cocktailID)
		if err != nil {
			return err
		}

		err = c.cocktailStepMySQLRepo.DeleteByCocktailIDTx(ctx, tx, cocktailID)
		if err != nil {
			return err
		}

		err = c.cocktailPhotoMySQLRepo.DeleteByCocktailIDTx(ctx, tx, cocktailID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *cocktailUsecase) MakeDraftToFormal(ctx context.Context, cocktailID, userID int64) error {

	cocktail, err := c.cocktailMySQLRepo.QueryByCocktailID(ctx, cocktailID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrCocktailNotFound
	} else if err != nil {
		return err
	}

	if cocktail.UserID != userID {
		return domain.ErrItemDoesNotBelongToUser
	}

	if cocktail.Category != cockarticletype.Draft.Int() {
		return domain.ErrCocktailNotFound
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
		return err
	}

	if apiCocktail.Title == "" || apiCocktail.Description == "" || len(apiCocktail.Ingredients) <= 0 ||
		len(apiCocktail.Steps) <= 0 || len(apiCocktail.Photos) <= 0 || len(apiCocktail.Photos) > 5 {
		return domain.ErrorCocktailNotFinished
	}

	user, err := c.userMySQLRepo.QueryById(ctx, userID)
	if err != nil {
		return err
	}

	if err := c.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		_, err := c.cocktailMySQLRepo.UpdateCategoryTx(ctx, tx,
			&domain.Cocktail{
				CocktailID: cocktail.CocktailID,
				Category:   cockarticletype.Normal.Int(),
			})
		if err != nil {
			return err
		}

		numberOfDraft := user.NumberOfDraft - 1
		_, err = c.userMySQLRepo.UpdateNumberOfDraftTx(ctx, tx,
			&domain.User{
				ID:            user.ID,
				NumberOfDraft: numberOfDraft,
			})
		if err != nil {
			return err
		}

		numberOfPost := user.NumberOfPost + 1
		_, err = c.userMySQLRepo.UpdateNumberOfPostTx(ctx, tx,
			&domain.User{
				ID:           user.ID,
				NumberOfPost: numberOfPost,
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
