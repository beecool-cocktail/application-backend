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

	_, _, err := util.DecodeBase64AndSaveAsWebp(ci.Data, "/"+ci.Destination)
	if err != nil {
		return err
	}

	return nil
}

func (u *cocktailFileRepository) SaveAsWebpInLQIP(ctx context.Context, ci *domain.CocktailImage) error {

	err := util.DecodeBase64AndSaveAsWebpInLQIP(ci.Data, "/"+ci.Destination)
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

func (u *cocktailFileRepository) UpdateAsWebpInLQIP(ctx context.Context, ci *domain.CocktailImage) error {

	err := util.DecodeBase64AndUpdateAsWebpInLQIP(ci.Data, "/"+ci.Destination)
	if err != nil {
		return err
	}

	return nil
}
