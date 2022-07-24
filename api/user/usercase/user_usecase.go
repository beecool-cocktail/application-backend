package usercase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userUsecase struct {
	service         *service.Service
	userMySQLRepo   domain.UserMySQLRepository
	userRedisRepo   domain.UserRedisRepository
	userFileRepo    domain.UserFileRepository
	transactionRepo domain.DBTransactionRepository
}

func NewUserUsecase(s *service.Service, clientMySQLRepo domain.UserMySQLRepository, clientRedisRepo domain.UserRedisRepository,
	userFileRepo domain.UserFileRepository, transaction domain.DBTransactionRepository) domain.UserUsecase {
	return &userUsecase{
		service:         s,
		userMySQLRepo:   clientMySQLRepo,
		userRedisRepo:   clientRedisRepo,
		userFileRepo:    userFileRepo,
		transactionRepo: transaction,
	}
}

func (u *userUsecase) Logout(ctx context.Context, id int64) (err error) {

	token := util.GenString(64)
	redisToken := domain.UserCache{
		Id:          id,
		AccessToken: token,
	}
	if err := u.userRedisRepo.UpdateToken(ctx, &redisToken); err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) QueryById(ctx context.Context, id int64) (domain.User, error) {

	user, err := u.userMySQLRepo.QueryById(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, domain.ErrUserNotFound
	} else if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *userUsecase) UpdateUserInfo(ctx context.Context, d *domain.User, ui *domain.UserImage) error {

	newFileName := uuid.New().String()

	savePath := u.service.Configure.Others.File.Image.PathInDB
	urlPath := u.service.Configure.Others.File.Image.PathInURL

	ui.Name = newFileName

	if !util.ValidateImageType(ui.Type) {
		return domain.ErrCodeFileTypeIllegal
	}

	err := u.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		if ui.Data != "" {
			ui.Destination = savePath + newFileName
			err := u.userFileRepo.SaveAsWebp(ctx, ui)
			if err != nil {
				return err
			}

			ui.Destination = urlPath + newFileName + ".webp"
			_, err = u.userMySQLRepo.UpdateImageTx(ctx, tx, ui)
			if err != nil {
				return err
			}
		}

		_, err := u.userMySQLRepo.UpdateBasicInfoTx(ctx, tx, d)
		if err != nil {
			return err
		}

		err = u.userRedisRepo.UpdateBasicInfo(ctx, &domain.UserCache{
			Name: d.Name,
		})
		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		return err
	}

	return nil
}
