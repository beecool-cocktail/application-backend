package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/cockarticletype"
	"github.com/beecool-cocktail/application-backend/enum/httpaction"
	"github.com/beecool-cocktail/application-backend/enum/sortbydir"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cocktailUsecase struct {
	service                     *service.Service
	cocktailMySQLRepo           domain.CocktailMySQLRepository
	cocktailRedisRepo           domain.CocktailRedisRepository
	cocktailElasticSearchRepo   domain.CocktailElasticSearchRepository
	cocktailFileRepo            domain.CocktailFileRepository
	cocktailPhotoMySQLRepo      domain.CocktailPhotoMySQLRepository
	cocktailIngredientMySQLRepo domain.CocktailIngredientMySQLRepository
	cocktailStepMySQLRepo       domain.CocktailStepMySQLRepository
	userMySQLRepo               domain.UserMySQLRepository
	favoriteCocktailMySQL       domain.FavoriteCocktailMySQLRepository
	transactionRepo             domain.DBTransactionRepository
}

// NewDietUsecase ...
func NewCocktailUsecase(
	s *service.Service,
	cocktailMySQLRepo domain.CocktailMySQLRepository,
	cocktailRedisRepo domain.CocktailRedisRepository,
	cocktailElasticSearchRepo domain.CocktailElasticSearchRepository,
	cocktailFileRepo domain.CocktailFileRepository,
	cocktailPhotoMySQLRepo domain.CocktailPhotoMySQLRepository,
	cocktailIngredientMySQLRepo domain.CocktailIngredientMySQLRepository,
	cocktailStepMySQLRepo domain.CocktailStepMySQLRepository,
	userMySQLRepo domain.UserMySQLRepository,
	favoriteCocktailMySQL domain.FavoriteCocktailMySQLRepository,
	transactionRepo domain.DBTransactionRepository) domain.CocktailUsecase {
	return &cocktailUsecase{
		service:                     s,
		cocktailMySQLRepo:           cocktailMySQLRepo,
		cocktailRedisRepo:           cocktailRedisRepo,
		cocktailElasticSearchRepo:   cocktailElasticSearchRepo,
		cocktailFileRepo:            cocktailFileRepo,
		cocktailPhotoMySQLRepo:      cocktailPhotoMySQLRepo,
		cocktailIngredientMySQLRepo: cocktailIngredientMySQLRepo,
		cocktailStepMySQLRepo:       cocktailStepMySQLRepo,
		userMySQLRepo:               userMySQLRepo,
		favoriteCocktailMySQL:       favoriteCocktailMySQL,
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

		lowQualityPhotos, err := c.cocktailPhotoMySQLRepo.QueryLowQualityPhotosByCocktailId(ctx, cocktail.CocktailID)
		if err != nil {
			return []domain.APICocktail{}, err
		}
		cocktail.LowQualityPhotos = lowQualityPhotos

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

func (c *cocktailUsecase) fillSelfCocktailList(ctx context.Context, cocktails []domain.APICocktail) ([]domain.APICocktail, error) {

	var apiCocktails []domain.APICocktail

	for _, cocktail := range cocktails {
		photos, err := c.cocktailPhotoMySQLRepo.QueryPhotosByCocktailId(ctx, cocktail.CocktailID)
		if err != nil {
			return []domain.APICocktail{}, err
		}
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

	lowQualityPhotos, err := c.cocktailPhotoMySQLRepo.QueryLowQualityPhotosByCocktailId(ctx, cocktail.CocktailID)
		if err != nil {
			return domain.APICocktail{}, err
		}
		cocktail.LowQualityPhotos = lowQualityPhotos

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

func (c *cocktailUsecase) fillCollectionStatusInDetails(ctx context.Context, cocktail domain.APICocktail, userID int64) (domain.APICocktail, error) {

	favoriteCocktails, _, err := c.favoriteCocktailMySQL.QueryByUserID(ctx, userID, domain.PaginationMySQLRepository{})
	if err != nil {
		return domain.APICocktail{}, err
	}

	for _, favoriteCocktail := range favoriteCocktails {
		if cocktail.CocktailID == favoriteCocktail.CocktailID {
			cocktail.IsCollected = true
		}
	}

	return cocktail, nil
}

func (c *cocktailUsecase) fillCollectionStatusInList(ctx context.Context, cocktails []domain.APICocktail,
	userID int64) ([]domain.APICocktail, error) {

	var apiCocktails []domain.APICocktail
	favoriteCocktails, _, err := c.favoriteCocktailMySQL.QueryByUserID(ctx, userID, domain.PaginationMySQLRepository{})
	if err != nil {
		return []domain.APICocktail{}, err
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

func (c *cocktailUsecase) getNeedDeletedPhoto(oldPhotos []domain.CocktailPhoto,
	newImages []domain.CocktailImage) []domain.CocktailPhoto {

	var deletedPhotoID []domain.CocktailPhoto
	for _, oldPhoto := range oldPhotos {
		var needDeleted = true
		for _, newImage := range newImages {
			if oldPhoto.ID == newImage.ImageID {
				needDeleted = false
			}
		}
		if needDeleted {
			deletedPhotoID = append(deletedPhotoID, oldPhoto)
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
	pathInServer := c.service.Configure.Others.File.Image.PathInServer
	pathInURL := c.service.Configure.Others.File.Image.PathInURL

	newFileName := uuid.New().String()
	lowQualityBundleID := util.GetID(util.IdGenerator)

	if !util.ValidateImageType(image.ContentType) {
		return domain.ErrCodeFileTypeIllegal
	}

	serverPath := util.ConcatString(pathInServer, newFileName)
	err := c.cocktailFileRepo.SaveAsWebp(ctx, 
		&domain.CocktailImage{
			File: image.File,
			ContentType: image.ContentType,
			Path: serverPath,
	})
	if err != nil {
		return err
	}

	imageType := util.GetImageType(image.ContentType)
	urlPath := util.ConcatString(pathInURL, newFileName, ".", imageType)
	err = c.cocktailPhotoMySQLRepo.StoreTx(ctx, tx,
		&domain.CocktailPhoto{
			CocktailID:         image.CocktailID,
			Photo:              urlPath,
			IsCoverPhoto:       image.IsCoverPhoto,
			IsLowQuality:       false,
			LowQualityBundleID: lowQualityBundleID,
			Order:              image.Order,
		})
	if err != nil {
		return err
	}

	lowQualityServerPath := util.ConcatString(pathInServer, newFileName)
	err = c.cocktailFileRepo.SaveAsWebpInLQIP(ctx, 
		&domain.CocktailImage{
			File: image.File,
			ContentType: image.ContentType,
			Path: lowQualityServerPath,
	})
	if err != nil {
		return err
	}

	lowQualityURLPath := util.ConcatString(pathInURL, newFileName, "_lq.", imageType)
	err = c.cocktailPhotoMySQLRepo.StoreTx(ctx, tx,
		&domain.CocktailPhoto{
			CocktailID:         image.CocktailID,
			Photo:              lowQualityURLPath,
			IsCoverPhoto:       image.IsCoverPhoto,
			IsLowQuality:       true,
			LowQualityBundleID: lowQualityBundleID,
			Order:              image.Order,
		})
	if err != nil {
		return err
	}

	return nil
}

func (c *cocktailUsecase) editPhoto(ctx context.Context, tx *gorm.DB, image *domain.CocktailImage) error {
	pathInServer := c.service.Configure.Others.File.Image.PathInServer
	pathInURL := c.service.Configure.Others.File.Image.PathInURL

	photo, err := c.cocktailPhotoMySQLRepo.QueryPhotoById(ctx, image.ImageID)
	if err != nil {
		return err
	}

	newFileName := uuid.New().String()
	oldFileName, err := util.GetFileNameByPath(photo.Photo)
	if err != nil {
		return err
	}

	imageType := util.GetImageType(image.ContentType)

	oldServerPath := util.ConcatString(pathInServer, oldFileName)
	newServerPath := util.ConcatString(pathInServer, newFileName, ".", imageType)
	err = c.cocktailFileRepo.UpdateAsWebp(ctx,
		&domain.CocktailImage{
			File:        image.File,
			ContentType: image.ContentType,
			Path:        oldServerPath,
		},
		newServerPath)
	if err != nil {
		return err
	}

	//this file already have type
	newURLPath := util.ConcatString(pathInURL, newFileName, ".", imageType)
	_, err = c.cocktailPhotoMySQLRepo.UpdateTx(ctx, tx,
		&domain.CocktailPhoto{
			ID:           photo.ID,
			Photo:        newURLPath,
			IsCoverPhoto: image.IsCoverPhoto,
			Order:        image.Order,
		})
	if err != nil {
		return err
	}

	lowQualityPhoto, err := c.cocktailPhotoMySQLRepo.QueryLowQualityPhotoByBundleId(ctx, photo.LowQualityBundleID)
	if err != nil {
		return err
	}

	oldLowQualityFileName, err := util.GetFileNameByPath(lowQualityPhoto.Photo)
	if err != nil {
		return err
	}

	lowQualityOldServerPath := util.ConcatString(pathInServer, oldLowQualityFileName)
	lowQualityNewServerPath := util.ConcatString(pathInServer, newFileName, "_lq.", imageType)
	err = c.cocktailFileRepo.UpdateAsWebpInLQIP(ctx, 
		&domain.CocktailImage{
		File:        image.File,
		ContentType: image.ContentType,
		Path:        lowQualityOldServerPath,
	},
	lowQualityNewServerPath)
	if err != nil {
		return err
	}

	newLowQualityURLPath := util.ConcatString(pathInURL, newFileName, "_lq.", imageType)
	_, err = c.cocktailPhotoMySQLRepo.UpdateTx(ctx, tx,
		&domain.CocktailPhoto{
			ID:           lowQualityPhoto.ID,
			Photo:        newLowQualityURLPath,
			IsCoverPhoto: image.IsCoverPhoto,
			Order:        image.Order,
		})
	if err != nil {
		return err
	}

	return nil
}

func (c *cocktailUsecase) GetAllWithFilter(ctx context.Context, filter map[string]interface{},
	pagination domain.PaginationUsecase, userID int64) ([]domain.APICocktail, int64, error) {

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

	if userID != 0 {
		apiCocktails, err = c.fillCollectionStatusInList(ctx, apiCocktails, userID)
		if err != nil {
			return []domain.APICocktail{}, 0, err
		}
	}

	return apiCocktails, total, nil
}

func (c *cocktailUsecase) Search(ctx context.Context, keyword string, from, size int, userID int64) ([]domain.APICocktail,
	int64, error) {

	cocktails, total, err := c.cocktailElasticSearchRepo.Search(ctx, keyword, from, size)
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

	if userID != 0 {
		apiCocktails, err = c.fillCollectionStatusInList(ctx, apiCocktails, userID)
		if err != nil {
			return []domain.APICocktail{}, 0, err
		}
	}

	return apiCocktails, total, nil
}

func (c *cocktailUsecase) QueryByCocktailID(ctx context.Context, cocktailID, userID int64) (domain.APICocktail, error) {

	cocktail, err := c.cocktailMySQLRepo.QueryByCocktailID(ctx, cocktailID)
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

	if userID != 0 {
		apiCocktail, err = c.fillCollectionStatusInDetails(ctx, apiCocktail, userID)
		if err != nil {
			return domain.APICocktail{}, err
		}
	}

	return apiCocktail, nil
}

func (c *cocktailUsecase) QueryFormalByUserID(ctx context.Context, targetUserID int64,
	queryUserID int64) ([]domain.APICocktail, error) {

	cocktails, err := c.cocktailMySQLRepo.QueryFormalByUserID(ctx, targetUserID)
	if err != nil {
		return []domain.APICocktail{}, err
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
		return []domain.APICocktail{}, err
	}

	if queryUserID != 0 {
		apiCocktails, err = c.fillCollectionStatusInList(ctx, apiCocktails, queryUserID)
		if err != nil {
			return []domain.APICocktail{}, err
		}
	}

	return apiCocktails, nil
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

func (c *cocktailUsecase) QueryFormalCountsByUserID(ctx context.Context, id int64) (int64, error) {

	total, err := c.cocktailMySQLRepo.QueryFormalCountsByUserID(ctx, id)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Todo 這裡不需要再傳userID
func (c *cocktailUsecase) Store(ctx context.Context, co *domain.Cocktail, ingredients []domain.CocktailIngredient,
	steps []domain.CocktailStep, images []domain.CocktailImage, userID int64) error {

	pathInServer := c.service.Configure.Others.File.Image.PathInServer
	pathInURL := c.service.Configure.Others.File.Image.PathInURL

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

		for order, image := range images {

			newFileName := uuid.New().String()
			image.Name = newFileName
			lowQualityBundleID := util.GetID(util.IdGenerator)

			if !util.ValidateImageType(image.ContentType) {
				return domain.ErrCodeFileTypeIllegal
			}

			image.Path = pathInServer + newFileName
			err := c.cocktailFileRepo.SaveAsWebp(ctx, &image)
			if err != nil {
				return err
			}

			imageType := util.GetImageType(image.ContentType)
			image.Path = util.ConcatString(pathInURL, newFileName, ".", imageType)
			err = c.cocktailPhotoMySQLRepo.StoreTx(ctx, tx,
				&domain.CocktailPhoto{
					CocktailID:         newCocktailID,
					Photo:              image.Path,
					IsCoverPhoto:       image.IsCoverPhoto,
					IsLowQuality:       false,
					LowQualityBundleID: lowQualityBundleID,
					Order:              order,
				})
			if err != nil {
				return err
			}

			image.Path = util.ConcatString(pathInServer, newFileName)
			err = c.cocktailFileRepo.SaveAsWebpInLQIP(ctx, &image)
			if err != nil {
				return err
			}

			image.Path = util.ConcatString(pathInURL, newFileName, "_lq.", imageType)
			err = c.cocktailPhotoMySQLRepo.StoreTx(ctx, tx,
				&domain.CocktailPhoto{
					CocktailID:         newCocktailID,
					Photo:              image.Path,
					IsCoverPhoto:       image.IsCoverPhoto,
					IsLowQuality:       true,
					LowQualityBundleID: lowQualityBundleID,
					Order:              order,
				})
			if err != nil {
				return err
			}
		}

		var elasticIngredients []string
		for _, ingredient := range ingredients {
			ingredient.CocktailID = newCocktailID
			err = c.cocktailIngredientMySQLRepo.StoreTx(ctx, tx, &ingredient)
			if err != nil {
				return err
			}
			elasticIngredients = append(elasticIngredients, ingredient.IngredientName)
		}

		var elasticSteps []string
		for _, step := range steps {
			step.CocktailID = newCocktailID
			err = c.cocktailStepMySQLRepo.StoreTx(ctx, tx, &step)
			if err != nil {
				return err
			}
			elasticSteps = append(elasticSteps, step.StepDescription)
		}

		if c.service.Configure.Elastic.Enable && co.Category == cockarticletype.Formal.Int() {
			err = c.cocktailElasticSearchRepo.Index(ctx, &domain.CocktailElasticSearch{
				CocktailID:  newCocktailID,
				UserID:      co.UserID,
				Title:       co.Title,
				Description: co.Description,
				Ingredients: elasticIngredients,
				Steps:       elasticSteps,
				CreatedDate: time.Now(),
			})
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

	oldPhoto, err := c.cocktailPhotoMySQLRepo.QueryPhotosByCocktailId(ctx, co.CocktailID)
	if err != nil {
		return err
	}

	deletedPhotos := c.getNeedDeletedPhoto(oldPhoto, images)

	if err := c.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		// update cocktail basic info
		_, err := c.cocktailMySQLRepo.UpdateTx(ctx, tx, co)
		if err != nil {
			return err
		}

		// update cocktail photo
		for _, image := range images {
			action, err := c.getAction(image.ImageID, image.File)
			if err != nil {
				return err
			}

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
				//keep photo and change order
				photo, err := c.cocktailPhotoMySQLRepo.QueryPhotoById(ctx, image.ImageID)
				if err != nil {
					return err
				}
				_, err = c.cocktailPhotoMySQLRepo.UpdatePhotoOrderTx(ctx, tx, &domain.CocktailPhoto{
					LowQualityBundleID: photo.LowQualityBundleID,
					Order:              image.Order,
					IsCoverPhoto:       image.IsCoverPhoto,
				})
				if err != nil {
					return err
				}
			}
		}

		for _, photo := range deletedPhotos {
			err := c.cocktailPhotoMySQLRepo.DeleteByLowQualityBundleIDTx(ctx, tx, photo.LowQualityBundleID)
			if err != nil {
				return err
			}
		}

		//update cocktail ingredient
		err = c.cocktailIngredientMySQLRepo.DeleteByCocktailIDTx(ctx, tx, co.CocktailID)
		if err != nil {
			return err
		}

		var elasticIngredients []string
		for _, ingredient := range ingredients {
			err = c.cocktailIngredientMySQLRepo.StoreTx(ctx, tx, &ingredient)
			if err != nil {
				return err
			}
			elasticIngredients = append(elasticIngredients, ingredient.IngredientName)
		}

		//update cocktail step
		err = c.cocktailStepMySQLRepo.DeleteByCocktailIDTx(ctx, tx, co.CocktailID)
		if err != nil {
			return err
		}

		var elasticSteps []string
		for _, step := range steps {
			err = c.cocktailStepMySQLRepo.StoreTx(ctx, tx, &step)
			if err != nil {
				return err
			}
			elasticSteps = append(elasticSteps, step.StepDescription)
		}

		if c.service.Configure.Elastic.Enable && co.Category == cockarticletype.Formal.Int() {
			err = c.cocktailElasticSearchRepo.Update(ctx, &domain.CocktailElasticSearch{
				CocktailID:  co.CocktailID,
				Title:       co.Title,
				Description: co.Description,
				Ingredients: elasticIngredients,
				Steps:       elasticSteps,
			})
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

		err = c.favoriteCocktailMySQL.DeleteTx(ctx, tx, cocktailID, domain.AllUsers)
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

		if c.service.Configure.Elastic.Enable && cocktail.Category == cockarticletype.Formal.Int(){
			err = c.cocktailElasticSearchRepo.Delete(ctx, cocktailID)
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
		len(apiCocktail.Steps) <= 0 || len(apiCocktail.Photos) > 5 {
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
				Category:   cockarticletype.Formal.Int(),
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

		if c.service.Configure.Elastic.Enable {
			err = c.cocktailElasticSearchRepo.Index(ctx, &domain.CocktailElasticSearch{
				CocktailID:  cocktail.CocktailID,
				UserID:      cocktail.UserID,
				Title:       cocktail.Title,
				Description: cocktail.Description,
				CreatedDate: time.Now(),
			})
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
