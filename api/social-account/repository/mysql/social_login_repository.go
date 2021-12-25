package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type socialAccountMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLSocialAccountRepository(db *gorm.DB) domain.SocialAccountMySQLRepository {
	return &socialAccountMySQLRepository{db}
}

func (s *socialAccountMySQLRepository) Store(ctx context.Context, ds *domain.SocialAccount, du *domain.User) (int64, error) {

	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Select("user_id", "account", "name", "email").Create(du).Error
		if err != nil {
			return err
		}

		ds.UserID = du.ID
		logrus.Debug(ds.UserID)
		err = tx.Select("social_id", "user_id", "type").Create(ds).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return du.ID, nil
}

func (s *socialAccountMySQLRepository) QueryById(ctx context.Context, id string) (*domain.SocialAccount, error) {
	var socialAccount domain.SocialAccount
	res := s.db.Where("social_id = ?", id).Take(&socialAccount)


	return &socialAccount, res.Error
}
