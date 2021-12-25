package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
)

type userMySQLRepository struct {
	db *gorm.DB
}

func NewMySQLUserRepository(db *gorm.DB) domain.UserMySQLRepository {
	return &userMySQLRepository{db}
}

func (u *userMySQLRepository) Store(ctx context.Context, d *domain.User) error {
	err := u.db.Select("id", "account", "password").Create(d).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *userMySQLRepository) QueryById(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	res := u.db.Where("id = ?", id).Find(&user)

	return &user, res.Error
}
