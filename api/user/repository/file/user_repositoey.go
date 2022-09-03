package file

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/util"
)

type userFileRepository struct {
}

func NewFileUserRepository() domain.UserFileRepository {
	return &userFileRepository{}
}

func (u *userFileRepository) SaveOriginAvatarAsWebp(ctx context.Context, o *domain.OriginAvatar) (int, int, error) {

	width, height, err := util.DecodeBase64AndSaveAsWebp(o.DataURL, "/"+o.Destination)
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func (u *userFileRepository) SaveCropAvatarAsWebp(ctx context.Context, c *domain.CropAvatar) error {

	_, _, err := util.DecodeBase64AndSaveAsWebp(c.DataURL, "/"+c.Destination)
	if err != nil {
		return err
	}

	return nil
}
