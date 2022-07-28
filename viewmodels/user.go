package viewmodels

type GoogleAuthenticateRequest struct {
	//required: true
	Code string `json:"code"`
}

type GoogleAuthenticateResponse struct {
	//required: true
	Token string `json:"token"`
}

type LogoutRequest struct {
	//required: true
	UserID int64 `json:"user_id"`
}

type GetUserInfoResponse struct {
	//required: true
	UserID int64 `json:"user_id"`
	//required: true
	Name string `json:"user_name"`
	//required: true
	Email string `json:"email"`
	//required: true
	Photo string `json:"photo"`
	//原圖長
	// required: true
	Length float32 `json:"length" binding:"required"`
	//原圖寬
	// required: true
	Width float32 `json:"width" binding:"required"`
	//座標 [左上XY, 右下XY]
	// required: true
	Coordinate []Coordinate `json:"coordinate" binding:"required,gte=2"`
	//required: true
	NumberOfPost int64 `json:"number_of_post"`
	//required: true
	NumberOfCollection int64 `json:"number_of_collection"`
	// 是否公開收藏 false=不公開, true=公開
	//required: true
	IsCollectionPublic bool `json:"is_collection_public"`
}

type GetOtherUserInfoResponse struct {
	//required: true
	UserID int64 `json:"user_id"`
	//required: true
	Name string `json:"user_name"`
	//required: true
	Photo string `json:"photo"`
	//原圖長
	// required: true
	Length float32 `json:"length" binding:"required"`
	//原圖寬
	// required: true
	Width float32 `json:"width" binding:"required"`
	//座標 [左上XY, 右下XY]
	// required: true
	Coordinate []Coordinate `json:"coordinate" binding:"required,gte=2"`
	//required: true
	NumberOfPost int64 `json:"number_of_post"`
	//required: true
	NumberOfCollection int64 `json:"number_of_collection"`
	// 是否公開收藏 false=不公開, true=公開
	//required: true
	IsCollectionPublic bool `json:"is_collection_public"`
}

type UpdateUserInfoRequest struct {
	File string `json:"file"`

	//原圖長
	// required: true
	Length float32 `json:"length" binding:"required"`

	//原圖寬
	// required: true
	Width float32 `json:"width" binding:"required"`

	//座標 [左上XY, 右下XY]
	// required: true
	Coordinate []Coordinate `json:"coordinate" binding:"required,gte=2"`

	// name for this user
	// required: true
	// example: Feen Lin
	Name string `json:"name" binding:"required"`

	// public user collection post
	// required: true
	// example: false
	IsCollectionPublic bool `json:"is_collection_public"`
}

type Coordinate struct {
	// required: true
	X float32 `json:"x"`
	// required: true
	Y float32 `json:"y"`
}

type UpdateUserInfoResponse struct {
	//required: true
	Photo string `json:"photo"`
}

type CollectArticleRequest struct {
	// cocktail id
	// required: true
	// example: 123456
	ID int64 `json:"id" binding:"required"`
}

type GetUserFavoriteCocktailListResponse struct {
	//required: true
	IsPublic bool `json:"is_public"`
	//required: true
	Total int64 `json:"total"`
	//required: true
	FavoriteCocktailList []FavoriteCocktail `json:"favorite_cocktail_list"`
}

type FavoriteCocktail struct {
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
}

type DeleteFavoriteCocktailResponse struct {
	//required: true
	CommandID string `json:"command_id"`
}
