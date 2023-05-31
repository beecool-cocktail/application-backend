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

func (u *userUsecase) UpdateUserAvatar(ctx context.Context, d *domain.User, ui *domain.UserAvatar) error {

	newFileName := uuid.New().String()

	savePath := u.service.Configure.Others.File.Image.PathInServer
	urlPath := u.service.Configure.Others.File.Image.PathInURL

	err := u.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		if ui.OriginAvatar.DataURL != "" {
			if !util.ValidateImageType(ui.OriginAvatar.Type) {
				return domain.ErrCodeFileTypeIllegal
			}

			ui.OriginAvatar.Destination = savePath + newFileName + "_origin"
			width, height, err := u.userFileRepo.SaveOriginAvatarAsWebp(ctx, &ui.OriginAvatar)
			if err != nil {
				return err
			}
			d.Width = width
			d.Height = height

			ui.OriginAvatar.Destination = urlPath + newFileName + "_origin." + util.GetImageType(ui.OriginAvatar.Type)
			_, err = u.userMySQLRepo.UpdateUserOriginAvatarTx(ctx, tx, ui)
			if err != nil {
				return err
			}
		}

		if !util.ValidateImageType(ui.CropAvatar.Type) {
			return domain.ErrCodeFileTypeIllegal
		}

		ui.CropAvatar.Destination = savePath + newFileName + "_crop"
		err := u.userFileRepo.SaveCropAvatarAsWebp(ctx, &ui.CropAvatar)
		if err != nil {
			return err
		}

		ui.CropAvatar.Destination = urlPath + newFileName + "_crop." + util.GetImageType(ui.CropAvatar.Type)
		_, err = u.userMySQLRepo.UpdateUserCropAvatarTx(ctx, tx, ui)
		if err != nil {
			return err
		}

		_, err = u.userMySQLRepo.UpdateUserAvatarInfoTx(ctx, tx, d)
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

func (u *userUsecase) DeleteUserAvatar(ctx context.Context, userID int64) error {
	err := u.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		_, err := u.userMySQLRepo.UpdateUserOriginAvatarTx(ctx, tx,
			&domain.UserAvatar{
				UserID: userID,
			})
		if err != nil {
			return err
		}

		_, err = u.userMySQLRepo.UpdateUserCropAvatarTx(ctx, tx,
			&domain.UserAvatar{
				UserID: userID,
			})
		if err != nil {
			return err
		}

		_, err = u.userMySQLRepo.UpdateUserAvatarInfoTx(ctx, tx,
			&domain.User{
				ID: userID,
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

func (u *userUsecase) UpdateUserName(ctx context.Context, d *domain.User) error {

	err := u.transactionRepo.Transaction(func(i interface{}) error {
		tx := i.(*gorm.DB)

		_, err := u.userMySQLRepo.UpdateUserNameTx(ctx, tx, d)
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

func (u *userUsecase) UpdateUserCollectionStatus(ctx context.Context, d *domain.User) error {

	_, err := u.userMySQLRepo.UpdateUserCollectionStatus(ctx, d)
	if err != nil {
		return err
	}

	return nil
}
