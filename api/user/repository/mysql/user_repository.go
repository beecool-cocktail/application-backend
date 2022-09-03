package mysql

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type name struct {
	Name string `structs:"name"`
}

type isCollectionPublic struct {
	IsCollectionPublic bool `structs:"is_collection_public"`
}

type photoInfo struct {
	Height       int     `structs:"height,omitempty"`
	Width        int     `structs:"width,omitempty"`
	CoordinateX1 float32 `structs:"coordinate_x1"`
	CoordinateY1 float32 `structs:"coordinate_y1"`
	CoordinateX2 float32 `structs:"coordinate_x2"`
	CoordinateY2 float32 `structs:"coordinate_y2"`
}

type originAvatar struct {
	OriginAvatar string `structs:"origin_avatar"`
}

type cropAvatar struct {
	CropAvatar string `structs:"crop_avatar"`
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

func (u *userMySQLRepository) QueryById(ctx context.Context, id int64) (domain.User, error) {
	var user domain.User
	res := u.db.Where("id = ?", id).Take(&user)

	return user, res.Error
}

func (u *userMySQLRepository) UpdateUserNameTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	var user domain.User
	updateColumn := name{
		Name: d.Name,
	}

	res := tx.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateUserCollectionStatus(ctx context.Context, d *domain.User) (int64, error) {
	var user domain.User
	updateColumn := isCollectionPublic{
		IsCollectionPublic: d.IsCollectionPublic,
	}

	res := u.db.Model(&user).Where("id = ?", d.ID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateUserOriginAvatarTx(ctx context.Context, tx *gorm.DB, d *domain.UserAvatar) (int64, error) {
	var user domain.User
	updateColumn := originAvatar{
		OriginAvatar: d.OriginAvatar.Destination,
	}

	res := tx.Model(&user).Where("id = ?", d.UserID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateUserCropAvatarTx(ctx context.Context, tx *gorm.DB, d *domain.UserAvatar) (int64, error) {
	var user domain.User
	updateColumn := cropAvatar{
		CropAvatar: d.CropAvatar.Destination,
	}

	res := tx.Model(&user).Where("id = ?", d.UserID).Updates(structs.Map(updateColumn))

	return res.RowsAffected, res.Error
}

func (u *userMySQLRepository) UpdateUserAvatarInfoTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	var user domain.User
	updateColumn := photoInfo{
		Height:       d.Height,
		Width:        d.Width,
		CoordinateX1: d.CoordinateX1,
		CoordinateY1: d.CoordinateY1,
		CoordinateX2: d.CoordinateX2,
		CoordinateY2: d.CoordinateY2,
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
