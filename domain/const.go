package domain

import "github.com/beecool-cocktail/application-backend/enum/sortbydir"

type PaginationUsecase struct {
	Page     int
	PageSize int
	SortByDir map[string]int
}

type PaginationMySQLRepository struct {
	Page     int
	PageSize int
	SortByDir map[string]sortbydir.SortByDir
}