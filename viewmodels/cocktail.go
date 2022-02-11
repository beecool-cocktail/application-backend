package viewmodels

type GetPopularCocktailListRequest struct {
	//required: true
	//example: 1
	Page int `json:"page"`

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

type PostArticleRequest struct {
	//required: true
	Files          []string             `json:"files" binding:"required,lte=5"`

	//required: true
	//example: Gin Tonic
	Name           string               `json:"name" binding:"required"`

	IngredientList []CocktailIngredient `json:"ingredient_list"`

	StepList       []CocktailStep       `json:"step_list"`

	//required: true
	//example: Very good to drink
	Description    string               `json:"description" binding:"required"`
}

type CocktailIngredient struct {
	//example: Gin Tonic
	Name   string  `json:"name"`

	//example: 1
	Amount float32 `json:"amount"`

	//example: cup
	Unit   string  `json:"unit"`
}

type CocktailStep struct {
	//example: 1
	Step        int    `json:"step"`

	//example: shake
	Description string `json:"description"`
}
