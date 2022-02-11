package usercase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userUsecase struct {
	userMySQLRepo   domain.UserMySQLRepository
	userRedisRepo   domain.UserRedisRepository
	userFileRepo    domain.UserFileRepository
	transactionRepo domain.DBTransactionRepository
}

func NewUserUsecase(clientMySQLRepo domain.UserMySQLRepository, clientRedisRepo domain.UserRedisRepository,
	userFileRepo domain.UserFileRepository, transaction domain.DBTransactionRepository) domain.UserUsecase {
	return &userUsecase{
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
		logrus.Error(err)
		return err
	}

	return nil
}

func (u *userUsecase) QueryById(ctx context.Context, id int64) (*domain.User, error) {

	user, err := u.userMySQLRepo.QueryById(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) UpdateUserInfo(ctx context.Context, d *domain.User, ui *domain.UserImage) error {

	newFileName := uuid.New().String()

	//Todo move to config
	savePath := "static/images/"
	urlPath := "static/"

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

			//Todo research why use pointer!!!!!!!!!!!!!!!!!!!!
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

		return err
	})

	if err != nil {
		return err
	}

	err = u.userRedisRepo.UpdateBasicInfo(ctx, &domain.UserCache{
		Name: d.Name,
	})
	if err != nil {
		return err
	}

	return nil
}
