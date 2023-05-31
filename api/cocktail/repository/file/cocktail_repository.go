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

	destination := util.ConcatString("/", ci.Path)
	_, _, err := util.DecodeBase64AndSaveAsWebp(ci.File, ci.ContentType, destination)
	if err != nil {
		return err
	}

	return nil
}

func (u *cocktailFileRepository) SaveAsWebpInLQIP(ctx context.Context, ci *domain.CocktailImage) error {

	destination := util.ConcatString("/", ci.Path)
	err := util.DecodeBase64AndSaveAsWebpInLQIP(ci.File, ci.ContentType, destination)
	if err != nil {
		return err
	}

	return nil
}

func (u *cocktailFileRepository) UpdateAsWebp(ctx context.Context, ci *domain.CocktailImage, destination string) error {

	source := util.ConcatString("/", ci.Path)
	destination = util.ConcatString("/", destination)
	err := util.DecodeBase64AndUpdateAsWebp(ci.File, ci.ContentType, source, destination)
	if err != nil {
		return err
	}

	return nil
}

func (u *cocktailFileRepository) UpdateAsWebpInLQIP(ctx context.Context, ci *domain.CocktailImage, destination string) error {

	source := util.ConcatString("/", ci.Path)
	destination = util.ConcatString("/", destination)
	err := util.DecodeBase64AndUpdateAsWebpInLQIP(ci.File, ci.ContentType, source, destination)
	if err != nil {
		return err
	}

	return nil
}
