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

func (u *userFileRepository) SaveAsWebp(ctx context.Context, ui *domain.UserImage) error {

	ui.Destination = ui.Destination + ".webp"
	err := util.SaveAsWebp(ui.Data, "/" + ui.Destination)
	if err != nil {
		return err
	}

	return nil
}