package domain

import "github.com/beecool-cocktail/application-backend/enum/sortbydir"

type PaginationUsecase struct {
	Page      int
	PageSize  int
	SortByDir map[string]int
}

type PaginationMySQLRepository struct {
	Page      int
	PageSize  int
	SortByDir map[string]sortbydir.SortByDir
}

const (
	AllUsers = -1
	NoUser   = 0

	// Elastic search weight
	TitleWeight       = 4
	IngredientWeight  = 3
	DescriptionWeight = 2
	StepWeight        = 1
)
