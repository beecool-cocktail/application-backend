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
	CocktailID     int64                `json:"cocktail_id"`
	UserID         int64                `json:"user_id"`
	UserName       string               `json:"user_name"`
	Title          string               `json:"title"`
	Photos         []string             `json:"photos"`
	IngredientList []CocktailIngredient `json:"ingredient_list"`
	CreatedDate    string               `json:"created_date"`
}

type GetDraftCocktailListRequest struct {
	//required: true
	//example: 1
	Page int `json:"page"`

	//required: true
	//example: 10
	PageSize int `json:"page_size"`
}

type GetDraftCocktailListResponse struct {
	Total             int64               `json:"total"`
	DraftCocktailList []DraftCocktailList `json:"draft_cocktail_list"`
}

type DraftCocktailList struct {
	CocktailID  int64  `json:"cocktail_id"`
	Photo       string `json:"photo"`
	Title       string `json:"title"`
	CreatedDate string `json:"created_date"`
}

type PostArticleRequest struct {
	Files []string `json:"files" binding:"lte=5"`

	//required: true
	//example: Gin Tonic
	Name string `json:"name" binding:"required"`

	IngredientList []CocktailIngredient `json:"ingredient_list"`

	StepList []CocktailStep `json:"step_list"`

	//example: Very good to drink
	Description string `json:"description"`
}

type PostDraftArticleRequest struct {
	Files []string `json:"files" binding:"lte=5"`

	//required: true
	//example: Gin Tonic
	Name string `json:"name" binding:"required"`

	IngredientList []CocktailIngredient `json:"ingredient_list"`

	StepList []CocktailStep `json:"step_list"`

	//example: Very good to drink
	Description string `json:"description"`
}

type CocktailIngredient struct {
	//example: Gin Tonic
	Name string `json:"name"`

	//example: 1
	Amount float32 `json:"amount"`

	//example: cup
	Unit string `json:"unit"`
}

type CocktailStep struct {
	//example: shake
	Description string `json:"description"`
}

type GetCocktailByIDRequest struct {
	// ID of an cocktail item
	//
	// In: path
	ID int64 `json:"id"`
}

type GetCocktailByIDResponse struct {
	CocktailID     int64                `json:"cocktail_id"`
	UserID         int64                `json:"user_id"`
	UserName       string               `json:"user_name"`
	Photos         []string             `json:"photos"`
	Title          string               `json:"title"`
	Description    string               `json:"description"`
	IngredientList []CocktailIngredient `json:"ingredient_list"`
	StepList       []CocktailStep       `json:"step_list"`
	CreatedDate    string               `json:"created_date"`
}

type GetCocktailDraftByIDRequest struct {
	// ID of an cocktail item
	//
	// In: path
	ID int64 `json:"id"`
}

type GetCocktailDraftByIDResponse struct {
	CocktailID     int64                `json:"cocktail_id"`
	Photos         []string             `json:"photos"`
	Title          string               `json:"title"`
	Description    string               `json:"description"`
	IngredientList []CocktailIngredient `json:"ingredient_list"`
	StepList       []CocktailStep       `json:"step_list"`
	CreatedDate    string               `json:"created_date"`
}
