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
	CocktailID     int64                         `json:"cocktail_id"`
	UserID         int64                         `json:"user_id"`
	UserName       string                        `json:"user_name"`
	Title          string                        `json:"title"`
	Photos         []string                      `json:"photos"`
	IngredientList []CocktailIngredientWithoutID `json:"ingredient_list"`
	IsCollected    bool                          `json:"is_collected"`
	CreatedDate    string                        `json:"created_date"`
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

	IngredientList []CocktailIngredientWithoutID `json:"ingredient_list"`

	StepList []CocktailStepWithoutID `json:"step_list"`

	//example: Very good to drink
	Description string `json:"description"`
}

type PostDraftArticleRequest struct {
	Files []string `json:"files" binding:"lte=5"`

	//example: Gin Tonic
	Name string `json:"name"`

	IngredientList []CocktailIngredientWithoutID `json:"ingredient_list"`

	StepList []CocktailStepWithoutID `json:"step_list"`

	//example: Very good to drink
	Description string `json:"description"`
}

type CocktailIngredientWithoutID struct {
	//example: Gin Tonic
	Name string `json:"name"`

	//example: 1 cup
	Amount string `json:"amount"`
}

type CocktailIngredientWithID struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

type CocktailStepWithoutID struct {
	//example: shake
	Description string `json:"description"`
}

type CocktailStepWithID struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type CocktailPhotoWithID struct {
	ID    int64  `json:"id"`
	Photo string `json:"path"`
}

type GetCocktailByIDRequest struct {
	// ID of an cocktail item
	//
	// In: path
	ID int64 `json:"id"`
}

type GetCocktailByIDResponse struct {
	CocktailID     int64                         `json:"cocktail_id"`
	UserID         int64                         `json:"user_id"`
	UserName       string                        `json:"user_name"`
	Photos         []CocktailPhotoWithID         `json:"photos"`
	Title          string                        `json:"title"`
	Description    string                        `json:"description"`
	IngredientList []CocktailIngredientWithoutID `json:"ingredient_list"`
	StepList       []CocktailStepWithoutID       `json:"step_list"`
	IsCollected    bool                          `json:"is_collected"`
	CreatedDate    string                        `json:"created_date"`
}

type GetCocktailDraftByIDRequest struct {
	// ID of an cocktail item
	//
	// In: path
	ID int64 `json:"id"`
}

type GetCocktailDraftByIDResponse struct {
	CocktailID     int64                         `json:"cocktail_id"`
	Photos         []CocktailPhotoWithID         `json:"photos"`
	Title          string                        `json:"title"`
	Description    string                        `json:"description"`
	IngredientList []CocktailIngredientWithoutID `json:"ingredient_list"`
	StepList       []CocktailStepWithoutID       `json:"step_list"`
	CreatedDate    string                        `json:"created_date"`
}

type UpdateDraftArticleRequest struct {
	Photos []CocktailPhotoWithID `json:"photos"`

	//example: Gin Tonic
	Name string `json:"name"`

	//required: true
	IngredientList []CocktailIngredientWithoutID `json:"ingredient_list"`

	//required: true
	StepList []CocktailStepWithoutID `json:"step_list"`

	//example: Very good to drink
	Description string `json:"description"`
}

type DeleteDraftArticleRequest struct {
	DeletedIds []int64 `json:"deleted_ids"`
}

type GetSelfCocktailListResponse struct {
	PopularCocktailList []SelfCocktailList `json:"popular_cocktail_list"`
}

type SelfCocktailList struct {
	CocktailID int64  `json:"cocktail_id"`
	UserName   string `json:"user_name"`
	Photo      string `json:"photo"`
	Title      string `json:"title"`
}
