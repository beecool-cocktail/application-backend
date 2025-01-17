package viewmodels

type GetPopularCocktailListRequest struct {
	//required: true
	//example: 1
	Page int `form:"page" json:"page"`

	//required: true
	//example: 10
	PageSize int `form:"page_size" json:"pageSize"`

	Keyword string `form:"keyword" json:"keyword"`
}

type GetPopularCocktailListResponse struct {
	//required: true
	Total int64 `json:"total"`
	//required: true
	PopularCocktailList []PopularCocktailList `json:"popular_cocktail_list"`
}

type PopularCocktailList struct {
	//required: true
	CocktailID int64 `json:"cocktail_id"`
	//required: true
	UserID int64 `json:"user_id"`
	//required: true
	UserName string `json:"user_name"`
	//required: true
	Title string `json:"title"`
	//required: true
	Photos []CocktailPhotoWithIDInResponse `json:"photos"`
	//required: true
	IngredientList []CocktailIngredientWithoutIDInResponse `json:"ingredient_list"`
	//required: true
	IsCollected bool `json:"is_collected"`
	//required: true
	CreatedDate string `json:"created_date"`
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
	//required: true
	Total int64 `json:"total"`
	//required: true
	DraftCocktailList []DraftCocktailList `json:"draft_cocktail_list"`
}

type DraftCocktailList struct {
	//required: true
	CocktailID int64 `json:"cocktail_id"`
	//required: true
	Photo string `json:"photo"`
	//required: true
	Title string `json:"title"`
	//required: true
	Description string `json:"description"`
	//required: true
	CreatedDate string `json:"created_date"`
}

type PostArticleRequest struct {
	Files []string `json:"files" binding:"lte=5"`

	//required: true
	//example: Gin Tonic
	Name string `json:"name" binding:"required"`

	IngredientList []CocktailIngredientWithoutIDInRequest `json:"ingredient_list"`

	StepList []CocktailStepWithoutIDInRequest `json:"step_list"`

	//example: Very good to drink
	Description string `json:"description"`
}

type PostDraftArticleRequest struct {
	Files []string `json:"files" binding:"lte=5"`

	//example: Gin Tonic
	Name string `json:"name"`

	IngredientList []CocktailIngredientWithoutIDInRequest `json:"ingredient_list"`

	StepList []CocktailStepWithoutIDInRequest `json:"step_list"`

	//example: Very good to drink
	Description string `json:"description"`
}

type MakeDraftArticleToFormalArticle struct {
	// ID of a cocktail item
	//
	// In: path
	ID int64 `uri:"id" binding:"required"`
}

type CocktailIngredientWithoutIDInRequest struct {
	//example: Gin Tonic
	//required: true
	Name string `json:"name"`

	//example: 1 cup
	//required: true
	Amount string `json:"amount"`
}

type CocktailIngredientWithoutIDInResponse struct {
	//example: Gin Tonic
	//required: true
	Name string `json:"name"`

	//example: 1 cup
	//required: true
	Amount string `json:"amount"`
}

type CocktailIngredientWithIDInRequest struct {
	//required: true
	ID int64 `json:"id"`
	//required: true
	Name string `json:"name"`
	//required: true
	Amount string `json:"amount"`
	//required: true
	Unit string `json:"unit"`
}

type CocktailIngredientWithIDInResponse struct {
	//required: true
	ID int64 `json:"id"`
	//required: true
	Name string `json:"name"`
	//required: true
	Amount string `json:"amount"`
	//required: true
	Unit string `json:"unit"`
}

type CocktailStepWithoutIDInRequest struct {
	//example: shake
	//required: true
	Description string `json:"description"`
}

type CocktailStepWithoutIDInResponse struct {
	//example: shake
	//required: true
	Description string `json:"description"`
}

type CocktailStepWithIDInRequest struct {
	//required: true
	ID int64 `json:"id"`
	//required: true
	Description string `json:"description"`
}

type CocktailStepWithIDInResponse struct {
	//required: true
	ID int64 `json:"id"`
	//required: true
	Description string `json:"description"`
}

type CocktailPhotoWithIDInRequest struct {
	ID        int64  `json:"id"`
	ImageFile string `json:"image_file"`
}

type CocktailPhotoWithIDInResponse struct {
	//required: true
	ID int64 `json:"id"`
	//required: true
	ImagePath string `json:"image_path"`
	//required: true
	BlurImageDataURL string `json:"blur_image_data_url"`
}

type DraftCocktailPhotoWithIDInResponse struct {
	//required: true
	ID int64 `json:"id"`
	//required: true
	ImagePath string `json:"image_path"`
}

type GetCocktailByIDRequest struct {
	// ID of a cocktail item
	//
	// In: path
	ID int64 `uri:"id" binding:"required"`
}

type GetCocktailByIDResponse struct {
	//required: true
	CocktailID int64 `json:"cocktail_id"`
	//required: true
	UserID int64 `json:"user_id"`
	//required: true
	UserName string `json:"user_name"`
	//required: true
	UserPhoto string `json:"user_photo"`
	//原圖長
	// required: true
	Height int `json:"height" binding:"required"`
	//原圖寬
	// required: true
	Width int `json:"width" binding:"required"`
	//座標 [左上XY, 右下XY]
	// required: true
	Coordinate []Coordinate `json:"coordinate" binding:"required,gte=2"`
	//旋轉角度
	// required: true
	Rotation float32 `json:"rotation" binding:"required"`
	//required: true
	Photos []CocktailPhotoWithIDInResponse `json:"photos"`
	//required: true
	Title string `json:"title"`
	//required: true
	Description string `json:"description"`
	//required: true
	IngredientList []CocktailIngredientWithoutIDInResponse `json:"ingredient_list"`
	//required: true
	StepList []CocktailStepWithoutIDInResponse `json:"step_list"`
	//required: true
	IsCollected bool `json:"is_collected"`
	//required: true
	CreatedDate string `json:"created_date"`
}

type GetCocktailDraftByIDRequest struct {
	// ID of a cocktail item
	//
	// In: path
	ID int64 `uri:"id" binding:"required"`
}

type GetCocktailDraftByIDResponse struct {
	//required: true
	CocktailID int64 `json:"cocktail_id"`
	//required: true
	Photos []DraftCocktailPhotoWithIDInResponse `json:"photos"`
	//required: true
	Title string `json:"title"`
	//required: true
	Description string `json:"description"`
	//required: true
	IngredientList []CocktailIngredientWithoutIDInResponse `json:"ingredient_list"`
	//required: true
	StepList []CocktailStepWithoutIDInResponse `json:"step_list"`
	//required: true
	CreatedDate string `json:"created_date"`
}

type UpdateDraftArticleUriRequest struct {
	// ID of a cocktail item
	//
	// In: path
	ID int64 `uri:"id" binding:"required"`
}

type UpdateDraftArticleRequest struct {
	Photos []CocktailPhotoWithIDInRequest `json:"photos"`

	//example: Gin Tonic
	Name string `json:"name"`

	//required: true
	IngredientList []CocktailIngredientWithoutIDInRequest `json:"ingredient_list"`

	//required: true
	StepList []CocktailStepWithoutIDInRequest `json:"step_list"`

	//example: Very good to drink
	Description string `json:"description"`
}

type DeleteDraftArticleRequest struct {
	DeletedIds []int64 `json:"deleted_ids"`
}

type DeleteFormalArticleRequest struct {
	DeletedIds []int64 `json:"deleted_ids"`
}

type UpdateFormalArticleUriRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type UpdateFormalArticleRequest struct {
	Photos []CocktailPhotoWithIDInRequest `json:"photos"`

	//example: Gin Tonic
	//required: true
	Name string `json:"name" binding:"required"`

	//required: true
	IngredientList []CocktailIngredientWithoutIDInRequest `json:"ingredient_list" binding:"required,gt=0"`

	//required: true
	StepList []CocktailStepWithoutIDInRequest `json:"step_list" binding:"required,gt=0"`

	//example: Very good to drink
	Description string `json:"description"`
}

type GetSelfCocktailListResponse struct {
	//required: true
	CocktailList []SelfCocktailList `json:"cocktail_list"`
}

type SelfCocktailList struct {
	//required: true
	CocktailID int64 `json:"cocktail_id"`
	//required: true
	UserName string `json:"user_name"`
	//required: true
	Photo string `json:"photo"`
	//required: true
	Title string `json:"title"`
	//required: true
	CreatedDate string `json:"created_date"`
}

type GetOtherCocktailListRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type GetOtherCocktailListResponse struct {
	//required: true
	CocktailList []OtherCocktailList `json:"cocktail_list"`
}

type OtherCocktailList struct {
	//required: true
	CocktailID int64 `json:"cocktail_id"`
	//required: true
	UserName string `json:"user_name"`
	//required: true
	Photo string `json:"photo"`
	//required: true
	Title string `json:"title"`
	//required: true
	IsCollected bool `json:"is_collected"`
	//required: true
	CreatedDate string `json:"created_date"`
}
