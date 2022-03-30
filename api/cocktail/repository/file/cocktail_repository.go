package file

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/util"
)

type cocktailFileRepository struct {
}

func NewFileUserRepository() domain.CocktailFileRepository {
	return &cocktailFileRepository{}
}

func (u *cocktailFileRepository) SaveAsWebp(ctx context.Context, ci *domain.CocktailImage) error {

	ci.Destination = ci.Destination + ".webp"
	err := util.DecodeBase64AndSaveAsWebp(ci.Data, "/"+ci.Destination)
	if err != nil {
		return err
	}

	return nil
}

func (u *cocktailFileRepository) UpdateAsWebp(ctx context.Context, ci *domain.CocktailImage) error {

	err := util.DecodeBase64AndUpdateAsWebp(ci.Data, "/"+ci.Destination)
	if err != nil {
		return err
	}

	return nil
}
