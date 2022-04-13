package viewmodels

type GoogleAuthenticateRequest struct {
	//required: true
	Code string `json:"code"`
}

type GoogleAuthenticateResponse struct {
	Token string `json:"token"`
}

type LogoutRequest struct {
	//required: true
	UserID int64 `json:"user_id"`
}

type GetUserInfoResponse struct {
	UserID             int64  `json:"user_id"`
	Name               string `json:"user_name"`
	Email              string `json:"email"`
	Photo              string `json:"photo"`
	NumberOfPost       int    `json:"number_of_post"`
	NumberOfCollection int    `json:"number_of_collection"`
	// 是否公開收藏 false=不公開, true=公開
	IsCollectionPublic bool `json:"is_collection_public"`
}

type UpdateUserInfoRequest struct {
	File string `json:"file"`

	// name for this user
	// required: true
	// example: Feen Lin
	Name string `json:"name" binding:"required"`

	// public user collection post
	// required: true
	// example: false
	IsCollectionPublic bool `json:"is_collection_public"`
}

type UpdateUserInfoResponse struct {
	Photo string `json:"photo"`
}

type CollectArticleRequest struct {
	// cocktail id
	// required: true
	// example: 123456
	ID int64 `json:"id" binding:"required"`
}

type GetUserFavoriteCocktailListResponse struct {
	Total                int64              `json:"total"`
	FavoriteCocktailList []FavoriteCocktail `json:"favorite_cocktail_list"`
}

type FavoriteCocktail struct {
	CocktailID int64  `json:"cocktail_id"`
	UserName   string `json:"user_name"`
	Photo      string `json:"photo"`
	Title      string `json:"title"`
}
