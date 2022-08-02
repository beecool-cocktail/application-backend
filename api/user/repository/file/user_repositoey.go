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

func (u *userFileRepository) SaveAsWebp(ctx context.Context, ui *domain.UserImage) (int, int, error) {

	width, height, err := util.DecodeBase64AndSaveAsWebp(ui.Data, "/"+ui.Destination)
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}
