package viewmodels

type GetPopularCocktailListRequest struct {
	//required: true
	//example: 1
	Page     int `json:"page"`

	//required: true
	//example: 10
	PageSize int `json:"page_size"`
}

type GetPopularCocktailListResponse struct {
	Total               int64                 `json:"total"`
	PopularCocktailList []PopularCocktailList `json:"popular_cocktail_list"`
}

type PopularCocktailList struct {
	CocktailID  int64  `json:"cocktail_id"`
	Photo       string `json:"photo"`
	Title       string `json:"title"`
	CreatedDate string `json:"created_date"`
}
