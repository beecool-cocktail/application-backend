package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type basicInfo struct {
	Name               string `structs:"name"`
	IsCollectionPublic bool   `structs:"is_collection_public"`
}

type photo struct {
	Photo string `structs:"photo"`
}

type postNumbers struct {
	NumberOfPost int `structs:"number_of_post"`
}

type draftNumbers struct {
	NumberOfDraft int `structs:"number_of_draft"`
}

type collectionsNumbers struct {
	NumberOfCollection int `structs:"number_of_collection"`
}

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
	res := u.db.Where("id = ?", id).Take(&user)

	return &user, res.Error
}

func (u *userMySQLRepository) UpdateBasicInfo(ctx context.Context, d *domain.User) (int64, error) {
	var user domain.User
	updateColumn := basicInfo{
		Name:               d.Name,
		IsCollectionPublic: d.IsCollectionPublic,
	}

	res := u.db.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateImage(ctx context.Context, d *domain.UserImage) (int64, error) {
	var user domain.User
	updateColumn := photo{
		Photo: d.Destination,
	}

	res := u.db.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateBasicInfoTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	var user domain.User
	updateColumn := basicInfo{
		Name:               d.Name,
		IsCollectionPublic: d.IsCollectionPublic,
	}

	res := tx.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateImageTx(ctx context.Context, tx *gorm.DB, d *domain.UserImage) (int64, error) {
	var user domain.User
	updateColumn := photo{
		Photo: d.Destination,
	}

	res := tx.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateNumberOfPostTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	var user domain.User
	updateColumn := postNumbers{
		NumberOfPost: d.NumberOfPost,
	}

	res := tx.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateNumberOfDraftTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	var user domain.User
	updateColumn := draftNumbers{
		NumberOfDraft: d.NumberOfDraft,
	}

	res := tx.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateNumberOfNumberOfCollectionTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	var user domain.User
	updateColumn := collectionsNumbers{
		NumberOfCollection: d.NumberOfCollection,
	}

	res := tx.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}
